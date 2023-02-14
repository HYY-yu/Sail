package svc_publish

import (
	"context"
	"fmt"
<<<<<<< HEAD
=======
	"github.com/HYY-yu/seckill.pkg/core"
>>>>>>> origin-feng/master
	"strconv"
	"strings"
	"time"

	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.pkg/pkg/encrypt"
	"github.com/HYY-yu/seckill.pkg/pkg/mysqlerr_helper"
	"github.com/gogf/gf/v2/errors/gerror"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/HYY-yu/sail/internal/service/sail/api/repo"
	"github.com/HYY-yu/sail/internal/service/sail/model"
	"github.com/HYY-yu/sail/internal/service/sail/storage"
)

// PublishSystem 发布系统
type PublishSystem interface {
	// EnterPublish
	// ConfigSystem 判断本次更新配置是否需要进入发布系统（判断条件：编辑的命名空间是否需要发布），进入发布系统则不走原来的配置编辑逻辑。
	EnterPublish(ctx context.Context, projectID, namespaceID, configID int, content string) error

	// QueryPublishConfig 查询配置的状态
	QueryPublishConfig(ctx context.Context, configID int) (model.PublishConfig, string, error)
}

// ConfigSystem 配置系统
type ConfigSystem interface {
	// ConfigEdit 配置变更回调，有历史记录
	// 做一个配置覆盖编辑，如果是回滚，则用发布前版本覆盖
	// 如果是全量发布，则用发布内容覆盖
	ConfigEdit()

	// ConfigKey 获取配置 key 格式
	ConfigKey(isPublic bool, projectGroupID int, projectKey string, namespaceName string, configName string, configType model.ConfigType) string
}

type PublishSvc struct {
	DB           db.Repo
	Store        storage.Repo
	configSystem ConfigSystem

	ConfigRepo        repo.ConfigRepo
	ProjectRepo       repo.ProjectRepo
	NamespaceRepo     repo.NamespaceRepo
	PublishRepo       repo.PublishRepo
	PublishConfigRepo repo.PublishConfigRepo
}

func NewPublishSvc(
	db db.Repo,
	store storage.Repo,
	cs ConfigSystem,
	cr repo.ConfigRepo,
	pr repo.ProjectRepo,
	nr repo.NamespaceRepo,
	pur repo.PublishRepo,
	puc repo.PublishConfigRepo,
) *PublishSvc {
	svc := &PublishSvc{
		DB:                db,
		Store:             store,
		configSystem:      cs,
		ConfigRepo:        cr,
		ProjectRepo:       pr,
		NamespaceRepo:     nr,
		PublishRepo:       pur,
		PublishConfigRepo: puc,
	}
	return svc
}

// EnterPublish 并发安全
// 进入发布，如果namespace尚未处于发布期，则自动进入发布期
// 将 config 加入 namespace 的发布期
// 如果 config 已加入，则更新 config 内容
// 如果 config 未加入，则加入
// 如果 config 已被锁定，则返回无法进入
func (p *PublishSvc) EnterPublish(ctx context.Context, projectID, namespaceID, configID int, content string) error {
	publishToken, err := p.initPublish(ctx, projectID, namespaceID)
	if err != nil {
		return gerror.Wrap(err, "initPublish")
	}

	puMgr := p.PublishRepo.Mgr(ctx, p.DB.GetDb())

	publish, err := puMgr.WithOptions(puMgr.WithPublishToken(publishToken)).Catch()
	if err != nil {
		// 过会儿再读一下，为写冲突的 goroutine 准备的逻辑
		for i := 0; i < 3; i++ {
			time.Sleep(10 * time.Millisecond)
			publish, err = puMgr.WithOptions(puMgr.WithPublishToken(publishToken)).Catch()
			if err == nil {
				break
			}
		}
		if err != nil {
			return err
		}
	}
	if publish.Status > model.PublishStatusRelease {
		return gerror.Newf("publish status wrong with publish_id %d", publish.ID)
	}

	tx := p.DB.GetDb().Begin()
	defer tx.Rollback()
	pucMgr := p.PublishConfigRepo.Mgr(ctx, tx)
	// BEGIN

	// SELECT * FROM publish_config WHERE publish_id = ? AND config_id = ?;
	publishConfig, err := pucMgr.WithOptions(pucMgr.WithPublishID(publish.ID), pucMgr.WithConfigID(configID)).
		Get()
	if err != nil {
		return err
	}
	project, namespace, err := p.getProjectAndNamespace(ctx, projectID, namespaceID)
	if err != nil {
		return err
	}
	cmgr := p.ConfigRepo.Mgr(ctx, tx)
	config, err := cmgr.WithOptions(cmgr.WithID(configID)).Catch()
	if err != nil {
		return err
	}
	configKey := p.configSystem.ConfigKey(
		config.IsPublic,
		config.ProjectGroupID,
		project.Key,
		namespace.Name,
		config.Name,
		model.ConfigType(config.ConfigType),
	)

	if publishConfig.ID == 0 {
		// INSERT with unique key
		gresp := p.Store.Get(ctx, configKey)
		if gresp.Err != nil {
			return err
		}
		if len(gresp.Value) == 0 {
			return gerror.New("not found key: " + configKey)
		}

		now := time.Now()
		publishConfig = model.PublishConfig{
			PublishID:          publish.ID,
			ConfigID:           configID,
			ConfigPreReversion: gresp.Revision,
			Status:             model.PublishStatusRelease,
			CreateTime:         now,
			UpdateTime:         now,
		}
		err := pucMgr.CreatePublishConfig(&publishConfig)
		if err != nil {
			// 已经插入则不管，继续流程
			if !mysqlerr_helper.IsMysqlDupEntryError(err) {
				return err
			}
		}
	} else {
		if publishConfig.Status > model.PublishStatusRelease {
			return gerror.Newf("publish status wrong with config_id %d", configID)
		}
	}

	// Update ETCD
	encryptContent := generatePublishContent(publishToken, publishConfig.ConfigPreReversion, content)

	sresp := p.Store.Set(ctx, configKey, encryptContent)
	if sresp.Err != nil {
		return sresp.Err
	}

	// Commit
	tx.Commit()
	return nil
}

func (p *PublishSvc) QueryPublishConfig(ctx context.Context, configID int) (model.PublishConfig, string, error) {
	return model.PublishConfig{}, "", nil
}

func (p *PublishSvc) LockPublish(ctx core.SvcContext, publishID int) {

}

func (p *PublishSvc) ListPublish() {

}

func (p *PublishSvc) RollbackPublish(ctx core.SvcContext) {

}

// initPublish 进入发布期
// 0. 检查 token
// 1. 生成 token
// 2. 写入 token
// 3. 写入发布表
// 幂等，可重入
func (p *PublishSvc) initPublish(ctx context.Context, projectID, namespaceID int) (string, error) {
	project, namespace, err := p.getProjectAndNamespace(ctx, projectID, namespaceID)
	if err != nil {
		return "", err
	}

	publishKey := publishTokenKey(project.Key, namespace.Name)

	gresp := p.Store.Get(ctx, publishKey)
	if gresp.Err != nil {
		return "", err
	}

	if len(gresp.Value) != 0 {
		return gresp.Value, nil
	}

	publishToken := generatePublishToken(projectID, namespaceID)

	sresp := p.Store.ConcurrentSet(ctx, publishKey, publishToken)
	if sresp.Err != nil {
		if sresp.Err != concurrency.ErrLocked {
			// 非写冲突直接返回
			return "", sresp.Err
		}
	}

	if sresp.Err == nil {
		puMgr := p.PublishRepo.Mgr(ctx, p.DB.GetDb())

		err := puMgr.CreatePublish(&model.Publish{
			ProjectID:    projectID,
			NamespaceID:  namespaceID,
			Status:       model.PublishStatusRelease,
			PublishToken: publishToken,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		})
		if err != nil {
			// 写入失败，删除这个 Token
			er := p.Store.Del(ctx, publishKey)
			if er != nil {
				return "", gerror.Wrap(err, "store err "+er.Error())
			}
			return "", err
		}
		return publishToken, nil
	}

	// 轮询读，为写冲突的 goroutine 准备的逻辑
	ta := time.After(3 * time.Second)
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			return "", gerror.New("Read key err: context canceled ")
		case <-ta:
			return "", gerror.New("Read timeout ")
		default:

		}
		gresp := p.Store.Get(ctx, publishKey)
		if len(gresp.Value) != 0 {
			return gresp.Value, nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return "", gerror.New("Read failed. ")
}

func (p *PublishSvc) getProjectAndNamespace(ctx context.Context, projectID int, namespaceID int) (*model.Project, *model.Namespace, error) {
	pMgr := p.ProjectRepo.Mgr(ctx, p.DB.GetDb())
	nMgr := p.NamespaceRepo.Mgr(ctx, p.DB.GetDb())
	pMgr.WithPrepareStmt()
	nMgr.WithPrepareStmt()

	project, err := pMgr.WithOptions(pMgr.WithID(projectID)).
		WithSelects(model.ProjectColumns.ID, model.ProjectColumns.Name).Catch()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "getConfigProjectAndNamespace")
	}
	namespace, err := nMgr.WithOptions(pMgr.WithID(namespaceID)).
		WithSelects(model.NamespaceColumns.ID, model.NamespaceColumns.Name).Catch()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "getConfigProjectAndNamespace")
	}

	return &project, &namespace, nil
}

// /conf/projectKey/namespaceName/publish/token
// 锁定期删除
func publishTokenKey(projectKey, namespaceName string) string {
	builder := strings.Builder{}
	builder.WriteString("/conf")

	builder.WriteByte('/')
	builder.WriteString(projectKey)

	builder.WriteByte('/')
	builder.WriteString(namespaceName)

	builder.WriteByte('/')
	builder.WriteString("publish")

	builder.WriteByte('/')
	builder.WriteString("token")

	return builder.String()
}

func generatePublishToken(projectID, namespaceID int) string {
	return encrypt.SHA256WithEncoding(fmt.Sprintf("%d-%d-%s", projectID, namespaceID, encrypt.Nonce(5)), encrypt.NewBase32Human())
}

// PUBLISH&publishToken{6}&pre-reversion&EncryptContent
func generatePublishContent(publishToken string, preReversion int, content string) string {
	builder := strings.Builder{}
	builder.WriteString("PUBLISH")
	builder.WriteByte('&')

	builder.WriteString(publishToken[:7])
	builder.WriteByte('&')
	builder.WriteString(strconv.Itoa(preReversion))

	goAES := encrypt.NewGoAES(publishToken, encrypt.AES192)
	encryptContent, err := goAES.WithModel(encrypt.ECB).WithEncoding(encrypt.NewBase64Encoding()).Encrypt(content)
	if err != nil {
		// 不可能为 err
		return ""
	}
	builder.WriteByte('&')
	builder.WriteString(encryptContent)

	return builder.String()
}

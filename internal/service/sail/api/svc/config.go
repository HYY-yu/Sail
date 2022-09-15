package svc

import (
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/db"

	"github.com/HYY-yu/sail/internal/service/sail/api/repo"
	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type ConfigSvc struct {
	BaseSvc
	DB db.Repo

	ConfigRepo        repo.ConfigRepo
	ConfigHistoryRepo repo.ConfigHistoryRepo
	ConfigLinkRepo    repo.ConfigLinkRepo
}

func NewConfigSvc(
	db db.Repo,
	cr repo.ConfigRepo,
	ch repo.ConfigHistoryRepo,
	cl repo.ConfigLinkRepo,
) *ConfigSvc {
	svc := &ConfigSvc{
		DB:                db,
		ConfigRepo:        cr,
		ConfigHistoryRepo: ch,
		ConfigLinkRepo:    cl,
	}
	return svc
}

// Add
// 1. 如果是公共配置，直接写到表里，把内容写到ETCD里
// 2. 如果加密，则需要在写内容之前，用 NamespaceID 的key去加密再写入
// 3. 如果依赖公共配置，则去公共配置中读取，并根据Link标识，决定是否插入到ConfigLink表
// 4. 如果不依赖公共配置，则同1
func (s *ConfigSvc) Add(sctx core.SvcContext, param *model.AddConfig) error {
	return nil
}

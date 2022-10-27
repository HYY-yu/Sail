package svc

import (
	"context"
	"net/http"
	"time"

	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.pkg/pkg/login"
	modellogin "github.com/HYY-yu/seckill.pkg/pkg/login/model"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"github.com/HYY-yu/seckill.pkg/pkg/token"
	"github.com/gogf/gf/v2/errors/gerror"
	"golang.org/x/crypto/bcrypt"

	"github.com/HYY-yu/sail/internal/service/sail/api/repo"
	"github.com/HYY-yu/sail/internal/service/sail/config"
	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type LoginSvc struct {
	DB        db.Repo
	StaffRepo repo.StaffRepo

	system login.LoginTokenSystem
}

func NewLoginSvc(
	db db.Repo,
	staffRepo repo.StaffRepo,
) *LoginSvc {
	svc := &LoginSvc{
		DB:        db,
		StaffRepo: staffRepo,
	}
	jwtCfg := config.Get().JWT

	switch jwtCfg.Type {
	case "refresh_token":
		cfg := &refreshTokenConfig{
			Secret:          jwtCfg.Secret,
			ExpireDuration:  jwtCfg.ExpireDuration,
			RefreshDuration: jwtCfg.RefreshDuration,
		}

		svc.system = NewByRefreshToken(cfg, staffRepo, db)
	}
	return svc
}

func (s *LoginSvc) Login(sctx core.SvcContext, param *model.LoginParams) (*model.LoginResponse, error) {
	ctx := sctx.Context()
	mgr := s.StaffRepo.Mgr(ctx, s.DB.GetDb())

	// 查询此用户
	user, err := mgr.WithOptions(mgr.WithName(param.UserName)).Get()
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if user.ID == 0 {
		return nil, response.NewErrorWithStatusOk(
			response.ServerError,
			"未找到用户",
		)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.Password))
	if err != nil {
		return nil, response.NewErrorWithStatusOk(
			response.ParamBindError,
			"输入的密码不正确",
		).WithErr(err)
	}

	// 派发Token
	resp, err := s.system.GenerateToken(ctx, user.ID, user.Name)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	lr := resp.Token.(*loginResponseByRefreshToken)
	result := &model.LoginResponse{
		AccessToken:  lr.AccessToken,
		RefreshToken: lr.RefreshToken,
		InitPassword: param.Password == "123456",
	}

	return result, err
}

func (s *LoginSvc) RefreshToken(sctx core.SvcContext, oldToken string) (*model.LoginResponse, error) {
	ctx := sctx.Context()
	// 检查refreshToken 是否过期、是否一致
	resp, err := s.system.RefreshToken(ctx, oldToken)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	lr := resp.Token.(*loginResponseByRefreshToken)
	result := &model.LoginResponse{
		AccessToken:  lr.AccessToken,
		RefreshToken: lr.RefreshToken,
	}

	return result, err
}

func (s *LoginSvc) LoginOut(sctx core.SvcContext) error {
	ctx := sctx.Context()
	userId := int(sctx.UserId())

	mgr := s.StaffRepo.Mgr(ctx, s.DB.GetDb())
	err := mgr.WithOptions(mgr.WithID(userId)).
		Update(model.StaffColumns.RefreshToken, "").Error
	if err != nil {
		return gerror.Wrap(err, "Update")
	}
	return nil
}

func (s *LoginSvc) ChangePassword(sctx core.SvcContext, newPass string) error {
	ctx := sctx.Context()
	mgr := s.StaffRepo.Mgr(ctx, s.DB.GetDb())
	if len(newPass) < 6 || len(newPass) > 10 {
		return response.NewErrorWithStatusOk(
			response.ParamBindError,
			"请输入6-10位密码 ",
		)
	}

	hp, err := bcrypt.GenerateFromPassword([]byte(newPass), 0)
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	err = mgr.WithOptions(mgr.WithID(int(sctx.UserId()))).
		Update(model.StaffColumns.Password, string(hp)).Error
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

type refreshTokenConfig struct {
	Secret          string        `json:"secret"`
	ExpireDuration  time.Duration `json:"expire_duration"`
	RefreshDuration time.Duration `json:"refresh_duration"`
}

type loginResponseByRefreshToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type loginSystemRefreshToken struct {
	cfg *refreshTokenConfig

	d db.Repo
	r repo.StaffRepo
}

func (l loginSystemRefreshToken) GenerateToken(ctx context.Context, userId int, userName string) (*modellogin.LoginResponse, error) {
	accessToken, err := token.New(l.cfg.Secret).JwtSign(int64(userId), userName, l.cfg.ExpireDuration)
	if err != nil {
		return nil, err
	}
	refreshToken, err := token.New(l.cfg.Secret).JwtSign(int64(userId), userName, l.cfg.RefreshDuration)
	if err != nil {
		return nil, err
	}

	// refreshToken 存到用户信息中
	mgr := l.r.Mgr(ctx, l.d.GetDb())
	err = mgr.WithOptions(mgr.WithID(userId)).
		Update(model.StaffColumns.RefreshToken, refreshToken).Error
	if err != nil {
		return nil, gerror.Wrap(err, "Update")
	}

	return &modellogin.LoginResponse{
		Token: &loginResponseByRefreshToken{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (l loginSystemRefreshToken) TokenCancelById(ctx context.Context, userId int, userName string) error {
	panic("implement me")
}

func (l loginSystemRefreshToken) TokenCancel(ctx context.Context, token string) error {
	panic("implement me")
}

func (l loginSystemRefreshToken) RefreshToken(ctx context.Context, oldToken string) (*modellogin.LoginResponse, error) {
	cliaims, err := token.New(l.cfg.Secret).JwtParse(oldToken)
	if err != nil {
		return nil, err
	}
	userId, userName := int(cliaims.UserID), cliaims.UserName

	// 检查数据库
	mgr := l.r.Mgr(ctx, l.d.GetDb())
	staff, _ := mgr.WithOptions(mgr.WithID(userId)).WithSelects(model.StaffColumns.ID, model.StaffColumns.RefreshToken).
		Get()
	if staff.RefreshToken != oldToken {
		return nil, gerror.New("token invalid. ")
	}

	return l.GenerateToken(ctx, userId, userName)
}

func NewByRefreshToken(cfg *refreshTokenConfig, r repo.StaffRepo, d db.Repo) login.LoginTokenSystem {
	return &loginSystemRefreshToken{cfg: cfg, r: r, d: d}
}

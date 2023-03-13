//go:build wireinject
// +build wireinject

//go:generate wire gen .
package api

import (
	"github.com/HYY-yu/sail/internal/service/sail/api/svc_interface"
	"github.com/HYY-yu/sail/internal/service/sail/api/svc_publish"
	"github.com/HYY-yu/seckill.pkg/cache"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/google/wire"

	"github.com/HYY-yu/sail/internal/service/sail/api/handler"
	"github.com/HYY-yu/sail/internal/service/sail/api/repo"
	"github.com/HYY-yu/sail/internal/service/sail/api/svc"
	"github.com/HYY-yu/sail/internal/service/sail/storage"
)

// initHandlers init Handlers.
func initHandlers(d db.Repo, c cache.Repo, store storage.Repo) (*Handlers, error) {
	panic(wire.Build(
		repo.NewProjectGroupRepo,
		repo.NewStaffRepo,
		repo.NewStaffGroupRelRepo,
		repo.NewProjectRepo,
		repo.NewNamespaceRepo,
		repo.NewConfigRepo,
		repo.NewConfigLinkRepo,
		repo.NewConfigHistoryRepo,
		repo.NewPublishRepo,
		repo.NewPublishConfigRepo,
		svc.NewProjectGroupSvc,
		svc.NewStaffSvc,
		svc.NewLoginSvc,
		svc.NewProjectSvc,
		svc.NewNamespaceSvc,
		svc.NewConfigSvc,
		svc_publish.NewPublishSvc,
		wire.Bind(new(svc_interface.PublishSystem), new(*svc_publish.PublishSvc)),
		wire.Bind(new(svc_interface.ConfigSystem), new(*svc.ConfigSvc)),
		handler.NewProjectGroupHandler,
		handler.NewStaffHandler,
		handler.NewLoginHandler,
		handler.NewProjectHandler,
		handler.NewNamespaceHandler,
		handler.NewConfigHandler,
		handler.NewIndexHandler,
		NewHandlers,
	))
}

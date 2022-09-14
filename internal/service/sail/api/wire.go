//go:build wireinject
// +build wireinject

//go:generate wire gen .
package api

import (
	"github.com/HYY-yu/seckill.pkg/cache"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/google/wire"

	"github.com/HYY-yu/sail/internal/service/sail/api/handler"
	"github.com/HYY-yu/sail/internal/service/sail/api/repo"
	"github.com/HYY-yu/sail/internal/service/sail/api/svc"
)

// initHandlers init Handlers.
func initHandlers(d db.Repo, c cache.Repo) (*Handlers, error) {
	panic(wire.Build(
		repo.NewProjectGroupRepo,
		repo.NewStaffRepo,
		repo.NewStaffGroupRelRepo,
		repo.NewProjectRepo,
		repo.NewNamespaceRepo,
		svc.NewProjectGroupSvc,
		svc.NewStaffSvc,
		svc.NewLoginSvc,
		svc.NewProjectSvc,
		svc.NewNamespaceSvc,
		handler.NewProjectGroupHandler,
		handler.NewStaffHandler,
		handler.NewLoginHandler,
		handler.NewProjectHandler,
		handler.NewNamespaceHandler,
		NewHandlers,
	))
}

//go:build wireinject
// +build wireinject

//go:generate wire gen .
package api

import (
	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/HYY-yu/seckill.pkg/cache"
	"github.com/HYY-yu/seckill.pkg/db"
)

// initHandlers init Handlers.
func initHandlers(l *zap.Logger, d db.Repo, c cache.Repo) (*Handlers, error) {
	panic(wire.Build(
		NewHandlers,
	))
}

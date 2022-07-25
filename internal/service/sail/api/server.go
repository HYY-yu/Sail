package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/HYY-yu/seckill.pkg/cache"
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/core/middleware"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.pkg/pkg/jaeger"
	"github.com/HYY-yu/seckill.pkg/pkg/metrics"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/HYY-yu/sail/internal/service/sail/config"
)

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

type Server struct {
	Logger      *zap.Logger
	HttpServer  *http.Server
	GrpcServer  *grpc.Server
	DB          db.Repo
	Cache       cache.Repo
	Trace       *trace.TracerProvider
	HTTPMiddles middleware.Middleware
}

func NewApiServer(logger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}
	s := &Server{}
	s.Logger = logger
	cfg := config.Get()

	dbRepo, err := db.New(&db.DBConfig{
		User:            cfg.MySQL.Base.User,
		Pass:            cfg.MySQL.Base.Pass,
		Addr:            cfg.MySQL.Base.Addr,
		Name:            cfg.MySQL.Base.Name,
		MaxOpenConn:     cfg.MySQL.Base.MaxOpenConn,
		MaxIdleConn:     cfg.MySQL.Base.MaxIdleConn,
		ConnMaxLifeTime: cfg.MySQL.Base.ConnMaxLifeTime,
		ServerName:      cfg.Server.ServerName,
	})
	if err != nil {
		logger.Fatal("new db err", zap.Error(err))
	}
	s.DB = dbRepo

	cacheRepo, err := cache.New(cfg.Server.ServerName, &cache.RedisConf{
		Addr:         cfg.Redis.Addr,
		Pass:         cfg.Redis.Pass,
		Db:           cfg.Redis.Db,
		MaxRetries:   cfg.Redis.MaxRetries,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})
	if err != nil {
		logger.Fatal("new cache err", zap.Error(err))
	}
	s.Cache = cacheRepo

	// Jaeger
	var tp *trace.TracerProvider
	if cfg.Jaeger.StdOut {
		tp, err = jaeger.InitStdOutForDevelopment(cfg.Server.ServerName, cfg.Jaeger.UdpEndpoint)
	} else {
		tp, err = jaeger.InitJaeger(cfg.Server.ServerName, cfg.Jaeger.UdpEndpoint)
	}
	if err != nil {
		logger.Error("jaeger error", zap.Error(err))
	}
	s.Trace = tp

	// Metrics
	sn := strings.Split(cfg.Server.ServerName, "_")
	metrics.InitMetrics(strings.Join(sn[:2], "_"))

	// Repo Svc Handler
	c, err := initHandlers(logger, s.DB, s.Cache)
	if err != nil {
		panic(err)
	}

	// HTTP Server
	opts := make([]core.Option, 0)
	opts = append(opts, core.WithEnableCors())
	opts = append(opts, core.WithRecordMetrics(metrics.RecordMetrics))
	if !cfg.Server.Pprof {
		opts = append(opts, core.WithDisablePProf())
	}
	engine, err := core.New(cfg.Server.ServerName, logger, opts...)
	if err != nil {
		panic(err)
	}
	// Init HTTP Middles
	s.HTTPMiddles = middleware.New(logger, cfg.JWT.Secret)

	// Route
	s.Route(c, engine)
	server := &http.Server{
		Handler: engine,
	}
	s.HttpServer = server

	return s, nil
}

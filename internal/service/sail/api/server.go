package api

import (
	"errors"
	"github.com/HYY-yu/sail/internal/service/sail/api/svc"
	"github.com/HYY-yu/sail/internal/service/sail/api/svc_interface"
	"html/template"
	"net/http"
	"strings"

	"github.com/HYY-yu/seckill.pkg/cache"
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/core/middleware"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.pkg/pkg/jaeger"
	"github.com/HYY-yu/seckill.pkg/pkg/metrics"
	"github.com/gin-gonic/gin"
	"github.com/vearutop/statigz"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/HYY-yu/sail/internal/service/sail/api/handler"
	"github.com/HYY-yu/sail/internal/service/sail/config"
	"github.com/HYY-yu/sail/internal/service/sail/storage"
	"github.com/HYY-yu/sail/ui"
)

type Handlers struct {
	projectGroupHandler *handler.ProjectGroupHandler
	staffHandler        *handler.StaffHandler
	loginHandler        *handler.LoginHandler
	projectHandler      *handler.ProjectHandler
	namespaceHandler    *handler.NamespaceHandler
	configHandler       *handler.ConfigHandler
	indexHandler        *handler.IndexHandler
}

func NewHandlers(
	projectGroupHandler *handler.ProjectGroupHandler,
	staffHandler *handler.StaffHandler,
	loginHandler *handler.LoginHandler,
	projectHandler *handler.ProjectHandler,
	namespaceHandler *handler.NamespaceHandler,
	configHandler *handler.ConfigHandler,
	indexHandler *handler.IndexHandler,
	publishSystem svc_interface.PublishSystem,
	configSvc *svc.ConfigSvc,
) *Handlers {
	configSvc.SetPublishSystem(publishSystem)
	return &Handlers{
		projectGroupHandler: projectGroupHandler,
		staffHandler:        staffHandler,
		loginHandler:        loginHandler,
		projectHandler:      projectHandler,
		namespaceHandler:    namespaceHandler,
		configHandler:       configHandler,
		indexHandler:        indexHandler,
	}
}

type Server struct {
	Logger      *zap.Logger
	HttpServer  *http.Server
	GrpcServer  *grpc.Server
	DB          db.Repo
	Cache       cache.Repo
	Storage     storage.Repo
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

	etcdRepo, err := storage.New(&storage.ETCDConfig{
		Endpoints:            cfg.ETCD.Endpoints,
		Username:             cfg.ETCD.Username,
		Password:             cfg.ETCD.Password,
		DialTimeout:          cfg.ETCD.DialTimeout,
		DialKeepAlive:        cfg.ETCD.DialKeepAlive,
		DialKeepAliveTimeout: cfg.ETCD.DialKeepAliveTimeout,
	})
	if err != nil {
		logger.Fatal("new etcd err", zap.Error(err))
	}
	s.Storage = etcdRepo

	//cacheRepo, err := cache.New(cfg.Server.ServerName, &cache.RedisConf{
	//	Addr:         cfg.Redis.Addr,
	//	Pass:         cfg.Redis.Pass,
	//	Db:           cfg.Redis.Db,
	//	MaxRetries:   cfg.Redis.MaxRetries,
	//	PoolSize:     cfg.Redis.PoolSize,
	//	MinIdleConns: cfg.Redis.MinIdleConn,
	//})
	//if err != nil {
	//	logger.Fatal("new cache err", zap.Error(err))
	//}
	//s.Cache = cacheRepo

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
	metrics.InitMetrics(cfg.Server.ServerName)

	// Repo Svc Handler
	c, err := initHandlers(s.DB, s.Cache, s.Storage)
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

	// HTTP Static Server
	staticEngine := gin.New()
	templateHTML, err := template.ParseFS(ui.TemplateFs, "template/**/**/*.html")
	if err != nil {
		panic(err)
	}
	staticEngine.SetHTMLTemplate(templateHTML)

	fServer := statigz.FileServer(ui.StaticFs)

	// Route
	s.Route(c, engine)
	s.RouteHTML(c, staticEngine)

	server := &http.Server{
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if strings.HasPrefix(request.URL.Path, "/static") {
				// 缓存控制+压缩控制
				//Serve existing static resource.
				fServer.ServeHTTP(writer, request)
				return
			} else if strings.HasPrefix(request.URL.Path, "/ui") {
				staticEngine.ServeHTTP(writer, request)
				return
			}
			engine.ServeHTTP(writer, request)
		}),
	}
	s.HttpServer = server

	return s, nil
}

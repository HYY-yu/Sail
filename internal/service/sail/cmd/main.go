package main

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/HYY-yu/sail/internal/service/sail/api"
	"github.com/HYY-yu/sail/internal/service/sail/config"

	"github.com/HYY-yu/seckill.pkg/pkg/logger"
	"github.com/HYY-yu/seckill.pkg/pkg/shutdown"
)

// @title                       配置中心
// @version                     1.0
// @description                 配置中心接口设计文档
// @contact.name                fengyu
// @contact.url                 https://hyy-yu.space
// @host                        localhost:8080
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
func main() {
	config.InitConfig()
	lp := findLogConfigOption()

	l, err := logger.New(lp...)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = l.Sync()
	}()

	// 初始化 API 服务
	s, err := api.NewApiServer(l)
	if err != nil {
		panic(err)
	}

	// 启动HTTP服务
	go func() {
		// 获取地址
		s.HttpServer.Addr = config.Get().Server.Host
		if err := s.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal("http server startup err", zap.Error(err))
		}
	}()

	// 优雅关闭
	shutdown.NewHook().Close(
		// 关闭 http server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := s.HttpServer.Shutdown(ctx); err != nil {
				l.Error("server shutdown err", zap.Error(err))
			}
		},

		// 关闭 grpc server
		func() {
			s.GrpcServer.GracefulStop()
		},

		// 关闭 db
		func() {
			if s.DB != nil {
				if err := s.DB.DbClose(); err != nil {
					l.Error("dbw close err", zap.Error(err))
				}
			}
		},

		// 关闭 cache
		func() {
			if s.Cache != nil {
				if err := s.Cache.Close(); err != nil {
					l.Error("cache close err", zap.Error(err))
				}
			}
		},
		// 关闭 Trace
		func() {
			if s.Trace != nil {
				// Do not make the application hang when it is shutdown.
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				if err := s.Trace.Shutdown(ctx); err != nil {
					l.Error("trace close err", zap.Error(err))
				}
			}
		},
	)
}

func findLogConfigOption() []logger.Option {
	C := config.Get()
	result := make([]logger.Option, 0)

	if !C.Log.Stdout {
		result = append(result, logger.WithDisableConsole())
	}

	if C.Log.JsonFormat {
		result = append(result, logger.WithJsonFormat())
	}

	switch C.Log.Level {
	case "DEBUG":
		result = append(result, logger.WithDebugLevel())
	case "INFO":
		result = append(result, logger.WithInfoLevel())
	case "WARN":
		result = append(result, logger.WithWarnLevel())
	case "ERROR":
		result = append(result, logger.WithErrorLevel())
	}

	result = append(result, logger.WithFileRotationP(C.Log.LogPath))
	return result
}

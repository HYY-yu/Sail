package config

import (
	"flag"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var config = new(Config)

type Config struct {
	MySQL struct {
		Base struct {
			MaxOpenConn     int
			MaxIdleConn     int
			ConnMaxLifeTime time.Duration
			Addr            string
			User            string
			Pass            string
			Name            string
		}
	}

	ETCD struct {
		Endpoints            []string
		Username             string
		Password             string
		DialTimeout          time.Duration
		DialKeepAlive        time.Duration
		DialKeepAliveTimeout time.Duration
	}

	JWT struct {
		Secret          string
		ExpireDuration  time.Duration
		Type            string
		RefreshDuration time.Duration
	}

	Log struct {
		LogPath    string
		Level      string
		Stdout     bool
		JsonFormat bool
	}

	Server struct {
		ServerName     string
		Host           string
		Pprof          bool
		HistoryListLen int64
	}

	Jaeger struct {
		UdpEndpoint string
		StdOut      bool
	}

	SDK struct {
		ConfigFilePath  string
		LogLevel        string
		MergeConfig     bool
		MergeConfigFile string
	}
}

var configPath string

func InitConfig() {
	flag.StringVar(&configPath, "config_path", ".", "path of cfg.toml")
	flag.Parse()

	viper.SetConfigName("cfg")
	viper.SetConfigType("toml")
	// 本地开发配置
	viper.AddConfigPath("./internal/service/sail/config")
	// 容器配置
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})
}

func Get() Config {
	return *config
}

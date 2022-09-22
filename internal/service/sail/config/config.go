package config

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var config = new(Config)

type Config struct {
	MySQL struct {
		Base struct {
			MaxOpenConn     int           `toml:"maxOpenConn"`
			MaxIdleConn     int           `toml:"maxIdleConn"`
			ConnMaxLifeTime time.Duration `toml:"connMaxLifeTime"`
			Addr            string        `toml:"addr"`
			User            string        `toml:"user"`
			Pass            string        `toml:"pass"`
			Name            string        `toml:"name"`
		} `toml:"base"`
	} `toml:"mysql"`

	Redis struct {
		Addr        string `toml:"addr"`
		Pass        string `toml:"pass"`
		Db          int    `toml:"db"`
		MaxRetries  int    `toml:"maxRetries"`
		PoolSize    int    `toml:"poolSize"`
		MinIdleConn int    `toml:"minIdleConn"`
	} `toml:"redis"`

	ETCD struct {
		Endpoints            []string      `toml:"endpoints"`
		Username             string        `toml:"username"`
		Password             string        `toml:"password"`
		DialTimeout          time.Duration `toml:"dialTimeout"`
		DialKeepAlive        time.Duration `toml:"dialKeepAlive"`
		DialKeepAliveTimeout time.Duration `toml:"dialKeepAliveTimeout"`
	}

	JWT struct {
		Secret          string        `toml:"secret"`
		ExpireDuration  time.Duration `toml:"expireDuration"`
		Type            string        `toml:"type"`
		RefreshDuration time.Duration `toml:"refreshDuration"`
	} `toml:"jwt"`

	Log struct {
		LogPath    string `toml:"logPath"`
		Level      string `toml:"level"`
		Stdout     bool   `toml:"stdout"`
		JsonFormat bool   `toml:"jsonFormat"`
	} `toml:"log"`

	Server struct {
		ServerName     string `toml:"serverName"`
		Host           string `toml:"host"`
		Pprof          bool   `toml:"pprof"`
		HistoryListLen int64  `toml:"historyListLen"`
	} `toml:"server"`

	Jaeger struct {
		UdpEndpoint string `toml:"udpEndpoint"`
		StdOut      bool   `toml:"stdOut"`
	} `toml:"jaeger"`
}

func InitConfig() {
	viper.SetConfigName("cfg")
	viper.SetConfigType("toml")
	// 本地开发配置
	viper.AddConfigPath("./internal/service/sail/config")
	// 容器配置
	viper.AddConfigPath(".")

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

package config

import (
	"fmt"
	"github.com/caarlos0/env/v7"
	"sync"
)

type Config struct {
	Listen struct {
		BindIP string `env:"TODO_BIND_IP" envDefault:"127.0.0.1"`
		Port   string `env:"TODO_PORT" envDefault:"7540"`
	}
	Fiber struct {
		Mode         string `env:"FIBER_MODE" envDefault:"debug"`
		BodyLimit    int    `env:"FIBER_BODY_LIMIT" envDefault:"99999999"`
		IdleTimeOut  int    `env:"FIBER_IDLE" envDefault:"20"`
		AllowOrigins string `env:"FIBER_ALLOW_ORIGINS" envDefault:"*"`
		AllowHeaders string `env:"FIBER_ALLOW_HEADERS" envDefault:""`
		AllowMethods string `env:"FIBER_ALLOW_METHODS" envDefault:"GET,POST,HEAD,PUT,DELETE,PATCH"`
	}
	Logger struct {
		Output string `env:"LOGGER_OUTPUT" envDefault:"json"`
	}
	DB struct {
		Path string `env:"TODO_DBFILE" envDefault:"./scheduler.db"`
	}
}

var s *Config
var once sync.Once

func Init() {
	once.Do(func() {
		s = &Config{}
		err := env.Parse(s)
		if err != nil {
			panic(fmt.Errorf("failed to parse config: %v", err))
		}
	})
}

func Get() *Config {
	if s == nil {
		Init()
	}
	return s
}

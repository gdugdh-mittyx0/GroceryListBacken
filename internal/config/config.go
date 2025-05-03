package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppPort       uint `env:"APP_PORT" env-default:"3000"`
	IsDevelopment bool `env:"IS_DEVELOPMENT" env-default:"false"`

	DatabaseHost     string `env:"DATABASE_HOST"`
	DatabaseDB       string `env:"DATABASE_DATABASE"`
	DatabasePassword string `env:"DATABASE_PASSWORD"`
	DatabaseUser     string `env:"DATABASE_USER"`
	DatabasePort     uint   `env:"DATABASE_PORT"`
}

func NewLoadConfig() (Config, error) {
	var cfg Config
	cleanenv.ReadConfig(".env", &cfg)
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, err
	}
	fmt.Println("config DATABASE_DATABASE:", cfg.DatabaseDB)
	return cfg, nil
}

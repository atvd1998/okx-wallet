package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		HTTPAddr string `env:"HTTP_ADDR" envDefault:""`
	}
}

func Load() (*Config, error) {
	var conf Config
	if err := env.Parse(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func MustLoadWithEnvFile(envFile string) (*Config, error) {
	if envFile != "" {
		godotenv.Load(envFile)
	}

	return Load()
}
func MustLoad() *Config {
	conf, err := Load()
	if err != nil {
		panic(err)
	}
	return conf
}

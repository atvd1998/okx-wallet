package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		HTTPAddr string `env:"HTTP_ADDR" envDefault:""`
	}
	OKX struct {
		APIKey     string `env:"OKX_API_KEY" envDefault:""`
		SecretKey  string `env:"OKX_SECRET_KEY" envDefault:""`
		AccessKey  string `env:"OKX_ACCESS_KEY" envDefault:""`
		PassPharse string `env:"OKX_PASSPHARSE" envDefault:""`
		Web3Url    string `env:"OKX_WEB3_URL" envDefault:""`
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

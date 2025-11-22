package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	StoragePath string `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env:"HTTP_ADDRESS" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := "./config/local.yml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist: %s", configPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error reading config: %s", err)
	}
	return &cfg
}

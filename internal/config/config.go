package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)


type Config struct {
	Env       string  `yaml:"env" env:"ENV" envDefault:"dev"`
	Version   string  `yaml:"version" env:"VERSION" end-default:"1.0.0"`
	AdminCode string  `yaml:"admin_code" env:"ADMIN_CODE" envRequired:"true"`
	TGToken   string `yaml:"tg_token" env:"TG_TOKEN" envRequired:"true"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		log.Fatalf("configuration file path is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("configuration file does not exist: %v", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read configuration file: %v", err)
	}
	return &cfg
}

func fetchConfigPath() string {
	var flagPath string
	flag.StringVar(&flagPath, "config", "", "path to config file")
	flag.Parse()

	switch {
	case os.Getenv("CONFIG_PATH") != "":
		return os.Getenv("CONFIG_PATH")
	case flagPath != "":
		return flagPath
	default:
		return "./configs/config.yaml"
	}
}

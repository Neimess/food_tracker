package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App        AppConfig      `yaml:"app"`
	Telegram   TelegramConfig `yaml:"telegram"`
	DB         DBConfig       `yaml:"db"`
	Cache      CacheConfig    `yaml:"cache"`
	HTTPServer HTTPServer     `yaml:"http_server"`
}

type AppConfig struct {
	Env     string `yaml:"env" env:"ENV" env-default:"dev"`
	Version string
}

type HTTPServer struct {
	Address         string        `yaml:"address" env:"HTTP_ADDRESS" env-default:"localhost:8080"`
	MaxHeaderBytes  int           `yaml:"max_header_bytes" env:"HTTP_MAX_HEADER_BYTES" env-default:"1048576"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-default:"10s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"5s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"  env:"HTTP_IDLE_TIMEOUT"   env-default:"60s"`
}

type TelegramConfig struct {
	Token        string  `yaml:"token" env:"TG_TOKEN" env-required:"true"`
	AllowedUsers []int64 `yaml:"allowed_users" env:"TG_ALLOWED_USERS"`
	WebAppURL    string  `yaml:"webapp_url" env:"WEBAPP_URL"`
	Address		 string  `yaml:"address" env:"TG_ADDRESS" env-default:"8090"`
	WebHookURL		 	 string  `yaml:"webhook_url" env:"TG_URL" env-required:"true"`
}

type DBConfig struct {
	Path string `yaml:"path" env:"DB_PATH" env-default:"storage.db"`
}

type CacheConfig struct {
	Path string `yaml:"path" env:"CACHE_PATH" env-default:".cache/food-tracker"`
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

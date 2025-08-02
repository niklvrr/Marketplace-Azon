package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	Env string `yaml:"env"`
}

type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type DatabaseConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Name           string `yaml:"name"`
	MaxConnections int    `yaml:"max_connections"`
}

type JWTConfig struct {
	Secret     string        `yaml:"secret"`
	Expiration time.Duration `yaml:"expiration"`
}

type LogConfig struct {
	LogLevel string `yaml:"log_level"`
	Format   string `yaml:"format"`
}

type Config struct {
	App      AppConfig      `yaml:"app"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"logging"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("configs/config.yaml")
	if err != nil {
		log.Fatal("Не удалось открыть config.yaml:", err)
	}
	defer file.Close()

	var cfg Config
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		log.Fatal("Не удалось декодировать YAML:", err)
	}

	if app_port := os.Getenv("APP_PORT"); app_port != "" {
		cfg.Server.Port, err = strconv.Atoi(app_port)
	}

	if db_host := os.Getenv("DB_HOST"); db_host != "" {
		cfg.Database.Host = db_host
	}

	if db_port := os.Getenv("DB_PORT"); db_port != "" {
		cfg.Database.Port, err = strconv.Atoi(db_port)
	}

	if db_username := os.Getenv("DB_USER"); db_username != "" {
		cfg.Database.User = db_username
	}

	if db_password := os.Getenv("DB_PASSWORD"); db_password != "" {
		cfg.Database.Password = db_password
	}

	if db_name := os.Getenv("DB_NAME"); db_name != "" {
		cfg.Database.Name = db_name
	}

	if jwt_secret := os.Getenv("JWT_SECRET"); jwt_secret != "" {
		cfg.JWT.Secret = jwt_secret
	}

	return &cfg, nil
}

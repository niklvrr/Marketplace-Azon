package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
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
	Url            string `yaml:"url"`
	Name           string `yaml:"name"`
	Password       string `yaml:"password"`
	User           string `yaml:"user"`
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

	if appHost := os.Getenv("APP_HOST"); appHost != "" {
		cfg.Server.Host = appHost
	}

	if appPort := os.Getenv("APP_PORT"); appPort != "" {
		cfg.Server.Port, err = strconv.Atoi(appPort)
		if err != nil {
			return nil, err
		}
	}

	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.Database.Host = dbHost
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		cfg.Database.Port, err = strconv.Atoi(dbPort)
		if err != nil {
			return nil, err
		}
	}

	if dbUrl := os.Getenv("DB_URL"); dbUrl != "" {
		cfg.Database.Url = dbUrl
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.Database.Name = dbName
	}

	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		cfg.Database.Password = dbPassword
	}

	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		cfg.Database.User = dbUser
	}

	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.JWT.Secret = jwtSecret
	}

	return &cfg, nil
}

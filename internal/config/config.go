package config

import (
	"github.com/spf13/viper"
	"time"
)

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Name           string `mapstructure:"name"`
	MaxConnections int    `mapstructure:"max_connections"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}

type LogConfig struct {
	LogLevel string `mapstructure:"log_level"`
	Format   string `mapstructure:"format"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"logging"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	//v.SetConfigType("env")
	//v.SetConfigFile(".env")
	//_ = v.ReadInConfig()

	v.AddConfigPath("configs")
	v.SetConfigName(".env")
	v.SetConfigType("env")

	v.AutomaticEnv()
	//v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	//v.SetConfigName("config")
	//v.SetConfigType("yaml")
	//v.AddConfigPath("configs/")
	//if err := v.ReadInConfig(); err != nil {
	//	return nil, err
	//}

	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Port                       string `mapstructure:"PORT"`
	DeleteConfigIfRecentlyUsed bool   `mapstructure:"DELETE_CONFIG_IF_RECENTLY_USED"`
	RecentUseDurationDays      int    `mapstructure:"RECENT_USE_DURATION_DAYS"`
}

type DatabaseConfig struct {
	Driver   string `mapstructure:"DB_DRIVER"`
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Dbname   string `mapstructure:"DB_NAME"`
}

type LoggerConfig struct {
	LogLevel string `mapstructure:"LOG_LEVEL"`
}

func GetConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var serverConfig ServerConfig
	var dbConfig DatabaseConfig
	var loggerConfig LoggerConfig
	if err := viper.Unmarshal(&serverConfig); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&dbConfig); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&loggerConfig); err != nil {
		return nil, err
	}
	cfg := &Config{
		Server:   serverConfig,
		Database: dbConfig,
		Logger:   loggerConfig,
	}

	return cfg, nil
}

func (c *Config) GetServerConfig() ServerConfig {
	return c.Server
}

func (c *Config) GetDatabaseConfig() DatabaseConfig {
	return c.Database
}

func (c *Config) GetLoggerConfig() LoggerConfig {
	return c.Logger
}

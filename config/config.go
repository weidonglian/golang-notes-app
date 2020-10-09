package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type PostgresConfig struct {
	Host           string `default:"" envconfig:"POSTGRES_HOST"`
	Port           int    `default:"5432" envconfig:"POSTGRES_PORT"`
	Username       string `default:"postgres" envconfig:"POSTGRES_USERNAME"`
	Password       string `default:"postgres" envconfig:"POSTGRES_PASSWORD"`
	DBName         string `default:"notes_app_dev" envconfig:"POSTGRES_DBNAME"`
	MigrationsPath string `default:"./migrations" envconfig:"POSTGRES_MIGRATIONS_PATH"`
}

func (c PostgresConfig) GetDataSourceName() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.Username, c.Password, c.DBName)
}

func (c PostgresConfig) IsValid() bool {
	return c.Host != "" && c.Port > 0
}

type Config struct {
	Postgres   PostgresConfig
	ServerPort string `default:"4000" envconfig:"SERVER_PORT"`
	RunEnv     string `default:"" envconfig:"RUN_ENV"`
}

var current *Config

func GetConfig() (Config, error) {
	if current != nil {
		return *current, nil
	}

	current = &Config{}
	if err := envconfig.Process("", current); err != nil {
		return *current, err
	}

	return *current, nil
}

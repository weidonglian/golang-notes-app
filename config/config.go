package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
	"path"
	"path/filepath"
)

type PostgresConfig struct {
	Host        string `default:"" envconfig:"POSTGRES_HOST"`
	Port        int    `default:"" envconfig:"POSTGRES_PORT"`
	Username    string `default:"" envconfig:"POSTGRES_USER"`
	Password    string `default:"" envconfig:"POSTGRES_PASSWORD"`
	DBName      string `default:"" envconfig:"POSTGRES_DB"`
	DataBaseURL string `default:"" envconfig:"DATABASE_URL"`
}

func (c PostgresConfig) GetDataSourceName() string {
	if c.DataBaseURL != "" {
		return c.DataBaseURL
	} else {
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.Username, c.Password, c.DBName)
	}

}

type Config struct {
	Postgres       PostgresConfig
	ServerPort     int    `default:"4000" envconfig:"PORT"`
	JWTSecret      string `default:"MaPoDouFu" envconfig:"JWT_SECRET"`
	MigrationsPath string `default:"./internal/db/migrations" envconfig:"MIGRATIONS_PATH"`
}

func DefaultTestConfig() *Config {
	// go test and ginkgo change the current directory to the test package directory
	// it will fail all the folder related handler. we need to figure out the root directory.
	rootDir := getRootProjectFolder()
	if rootDir == "" {
		panic("Failed to get the project root dir")
	}
	return &Config{
		Postgres: PostgresConfig{
			Host:        "localhost",
			Port:        5434,
			Username:    "postgres",
			Password:    "postgres",
			DBName:      "postgres",
			DataBaseURL: "",
		},
		ServerPort:     4040,
		JWTSecret:      "@Test@NoteApp",
		MigrationsPath: path.Join(rootDir, "./internal/db/migrations"),
	}
}

func DefaultConfig() *Config {
	return &Config{
		Postgres: PostgresConfig{
			Host:        "localhost",
			Port:        5433,
			Username:    "postgres",
			Password:    "postgres",
			DBName:      "postgres",
			DataBaseURL: "",
		},
		ServerPort:     4000,
		JWTSecret:      "@Dev@NoteApp",
		MigrationsPath: "./internal/db/migrations",
	}
}

func MustParseFromEnv(config *Config) {
	if err := envconfig.Process("", config); err != nil {
		panic(err)
	}
}

func getRootProjectFolder() string {
	workDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for i := 0; i < 10; i++ {
		if _, err := os.Stat(path.Join(workDir, "go.mod")); err == nil {
			return workDir
		} else {
			workDir = filepath.Dir(workDir)
		}
	}
	return ""
}

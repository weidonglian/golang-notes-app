package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"path"
	"path/filepath"
)

type appMode string

const (
	appModeDev  appMode = "Dev"
	appModeTest appMode = "Test"
	appModeProd appMode = "Prod"
)

// There are a couple of ways to adapt the RUN_ENV/APP_MODE.
// 1. Set the `APP_MODE` to candidates <Dev, Prod, Test>
// 2. By default it will be `Dev` mode.
// 3. For production, you have to set correctly all the environment variables in envconfig tags.
// 3. Call `SetTestMode` in the API to switch to test mode for testing or writing unit tests.
var currentAppMode appMode = appModeDev

func SetTestMode() {
	currentAppMode = appModeTest
	currentConfig = nil
}

func GetAppMode() string {
	return string(currentAppMode)
}

func IsDevMode() bool {
	switch currentAppMode {
	case appModeProd:
		return false
	case appModeTest:
		return false
	case appModeDev:
		return true
	default:
		return true
	}
}

func IsProdMode() bool {
	return currentAppMode == appModeProd
}

func IsTestMode() bool {
	return currentAppMode == appModeTest
}

func init() {
	appModeFromEnv := os.Getenv("APP_MODE")
	if appModeFromEnv != "" {
		currentAppMode = appMode(appModeFromEnv)
	}
}

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
	ServerPort     int    `default:"4000" envconfig:"SERVER_PORT"`
	JWTSecret      string `default:"MaPoDouFu" envconfig:"JWT_SECRET"`
	MigrationsPath string `default:"./db/migrations" envconfig:"MIGRATIONS_PATH"`
}

var (
	defaultTestConfig = Config{
		Postgres: PostgresConfig{
			Host:        "localhost",
			Port:        5433,
			Username:    "postgres",
			Password:    "postgres",
			DBName:      "postgres",
			DataBaseURL: "",
		},
		ServerPort:     3000,
		JWTSecret:      "@Test@NoteApp",
		MigrationsPath: "./db/migrations",
	}

	defaultDevConfig = Config{
		Postgres: PostgresConfig{
			Host:        "localhost",
			Port:        5432,
			Username:    "postgres",
			Password:    "postgres",
			DBName:      "postgres",
			DataBaseURL: "",
		},
		ServerPort:     4000,
		JWTSecret:      "@Dev@NoteApp",
		MigrationsPath: "./db/migrations",
	}
)

var currentConfig *Config

func GetConfig() Config {
	if currentConfig != nil {
		return *currentConfig
	}

	switch currentAppMode {
	case appModeTest:
		// go test and ginkgo change the current directory to the test package directory
		// it will fail all the folder related handler. we need to figure out the root directory.
		rootDir := getRootProjectFolder()
		if rootDir == "" {
			panic("Failed to get the project root dir")
		}
		currentConfig = &defaultTestConfig
		currentConfig.MigrationsPath = path.Join(rootDir, currentConfig.MigrationsPath)
	case appModeDev:
		currentConfig = &defaultDevConfig
	case appModeProd:
		currentConfig = &Config{}
		if err := envconfig.Process("", currentConfig); err != nil {
			log.Fatal("AppMode: can not process the envconfig with error:", err)
			panic(err)
		}
	default:
		panic("unknown app mod")
	}
	return *currentConfig
}

func getRootProjectFolder() string {
	workDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for i := 0; i < 10; {
		if _, err := os.Stat(path.Join(workDir, "go.mod")); err == nil {
			return workDir
		} else {
			workDir = filepath.Dir(workDir)
		}
	}
	return ""
}

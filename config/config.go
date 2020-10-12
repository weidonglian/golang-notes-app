package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
)

type appMode string

const (
	appModeDEV  appMode = "Dev"
	appModeTest appMode = "Test"
	appModeProd appMode = "Prod"
)

// There are a couple of ways to adapt the RUN_ENV/APP_MODE.
// 1. Set the `APP_MODE` to candidates <Dev, Prod, Test>
// 2. By default it will be `Dev` mode.
// 3. For production, you have to set correctly all the environment variables in envconfig tags.
// 3. Call `SetTestMode` in the API to switch to test mode for testing or writing unit tests.
var currentAppMode appMode = appModeDEV

func SetTestMode() {
	currentAppMode = appModeTest
	currentConfig = nil
}

func IsDevMode() bool {
	switch currentAppMode {
	case appModeProd:
		return false
	case appModeTest:
		return false
	case appModeDEV:
		return true
	default:
		return true
	}
}

func IsProdMode() bool {
	return currentAppMode == appModeDEV
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
	Host     string `default:"" envconfig:"POSTGRES_HOST"`
	Port     int    `default:"" envconfig:"POSTGRES_PORT"`
	Username string `default:"" envconfig:"POSTGRES_USERNAME"`
	Password string `default:"" envconfig:"POSTGRES_PASSWORD"`
	DBName   string `default:"" envconfig:"POSTGRES_DBNAME"`
}

type Sqlite3Config struct {
	SourceName string `default:":memory:" envconfig:"SQLITE3_SOURCE_NAME"`
}

func (c PostgresConfig) GetDataSourceName() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.Username, c.Password, c.DBName)
}

const (
	DatabaseDriverSqlite3  = "sqlite3"
	DatabaseDriverPostgres = "postgres"
)

type Config struct {
	Postgres       PostgresConfig
	Sqlite3        Sqlite3Config
	DatabaseDriver string `default:"" envconfig:"DB_DRIVER""`
	ServerPort     string `default:"4000" envconfig:"SERVER_PORT"`
	JWTSecret      string `default:"MaPoDouFu" envconfig:"JWT_SECRET"`
	MigrationsPath string `default:"./db/migrations" envconfig:"MIGRATIONS_PATH"`
}

var (
	defaultTestConfig = Config{
		Postgres: PostgresConfig{},
		Sqlite3: Sqlite3Config{
			SourceName: ":memory:",
		},
		DatabaseDriver: DatabaseDriverSqlite3,
		ServerPort:     "3000",
		JWTSecret:      "@Test@NoteApp",
		MigrationsPath: "./db/migrations",
	}

	defaultDevConfig = Config{
		Postgres: PostgresConfig{},
		Sqlite3: Sqlite3Config{
			SourceName: "notes_app_dev.db",
		},
		DatabaseDriver: DatabaseDriverSqlite3,
		ServerPort:     "4000",
		JWTSecret:      "@Dev@NoteApp",
		MigrationsPath: "./db/migrations",
	}
)

var currentConfig *Config

func GetConfig() (Config, error) {
	if currentConfig != nil {
		return *currentConfig, nil
	}

	switch currentAppMode {
	case appModeTest:
		currentConfig = &defaultTestConfig
	case appModeDEV:
		currentConfig = &defaultDevConfig
	case appModeProd:
		currentConfig = &Config{}
		if err := envconfig.Process("", currentConfig); err != nil {
			log.Fatal("AppMode: can not process the envconfig with error:", err)
			panic(err)
		}
	}
	return *currentConfig, nil
}

package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	databaseHost     = "APP_DB_HOST"
	databasePort     = "APP_DB_PORT"
	databaseName     = "APP_DB_NAME"
	databaseUser     = "APP_DB_USER"
	databasePassword = "APP_DB_PASSWORD"

	defaultDatabaseHost = "localhost"
	defaultDatabasePort = "5432"

	timeZone = "APP_TIMEZONE"

	defaultTimezone = "Europe/Moscow"

	appPort = "APP_PORT"
)

func Init() error {
	if err := godotenv.Overload(".env"); err != nil {
		return err
	}

	if _, err := os.Stat(".env.local"); err == nil {
		if err := godotenv.Overload(".env.local"); err != nil {
			return err
		}
	}

	return nil
}

func Timezone() string {
	timezone, ok := os.LookupEnv(timeZone)
	if !ok {
		timezone = defaultTimezone
	}

	return timezone
}

func DataSourceName() string {
	user, ok := os.LookupEnv(databaseUser)
	if !ok {
		panic(fmt.Errorf("env.DatabaseUser is required"))
	}

	password, ok := os.LookupEnv(databasePassword)
	if !ok {
		panic(fmt.Errorf("env.DatabasePassword is required"))
	}

	name, ok := os.LookupEnv(databaseName)
	if !ok {
		panic(fmt.Errorf("env.DatabaseName is required"))
	}

	port, ok := os.LookupEnv(databasePort)
	if !ok {
		port = defaultDatabasePort
	}

	host, ok := os.LookupEnv(databaseHost)
	if !ok {
		port = defaultDatabaseHost
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=%s",
		user,
		password,
		host,
		port,
		name,
		Timezone(),
	)
}

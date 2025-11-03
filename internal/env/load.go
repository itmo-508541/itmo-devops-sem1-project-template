package env

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	DatabaseHost     = "APP_DB_HOST"
	DatabasePort     = "APP_DB_PORT"
	DatabaseName     = "APP_DB_NAME"
	DatabaseUser     = "APP_DB_USER"
	DatabasePassword = "APP_DB_PASSWORD"
	PortEnv          = "APP_PORT"
	dotEnv           = ".env"
	dotEnvLocal      = ".env.local"
)

func Load() error {
	if err := godotenv.Overload(dotEnv); err != nil {
		return err
	}

	if _, err := os.Stat(dotEnvLocal); err == nil {
		if err := godotenv.Overload(dotEnvLocal); err != nil {
			return err
		}
	}

	return nil
}

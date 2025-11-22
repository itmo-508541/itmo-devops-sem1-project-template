package settings

import (
	"fmt"
	"project_sem/internal/config"
)

const (
	DatabaseHostDefault    = "localhost"
	DatabasePortDefault    = "5432"
	DatabaseSslModeDefault = "disable"

	TimezoneDefault = "Europe/Moscow"

	WebHostDefault = "0.0.0.0"
	WebPortDefault = "8080"

	databaseHostEnv     = "APP_DB_HOST"
	databasePortEnv     = "APP_DB_PORT"
	databaseSslModeEnv  = "APP_DB_SSL_MODE"
	databaseNameEnv     = "APP_DB_NAME"
	databaseUserEnv     = "APP_DB_USER"
	databasePasswordEnv = "APP_DB_PASSWORD"

	timezoneEnv = "APP_TIMEZONE"

	webHostEnv = "APP_HOST"
	webPortEnv = "APP_PORT"
)

func DatabaseSourceName() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=%s",
		config.RequiredEnv(databaseUserEnv),
		config.RequiredEnv(databasePasswordEnv),
		config.OptionalEnv(databaseHostEnv, DatabaseHostDefault),
		config.OptionalEnv(databasePortEnv, DatabasePortDefault),
		config.RequiredEnv(databaseNameEnv),
		config.OptionalEnv(databaseSslModeEnv, DatabaseSslModeDefault),
		config.OptionalEnv(timezoneEnv, TimezoneDefault),
	)
}

func WebServerAddr() string {
	return fmt.Sprintf(
		"%s:%s",
		config.OptionalEnv(webHostEnv, WebHostDefault),
		config.OptionalEnv(webPortEnv, WebPortDefault),
	)
}

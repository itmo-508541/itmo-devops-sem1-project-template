package database

import (
	"context"
	"project_sem/internal/config"
	db "project_sem/internal/database"
	"project_sem/internal/database/command/migrate"
	"project_sem/internal/models/price"
	"project_sem/internal/models/report"
	"project_sem/internal/services/general"

	"github.com/sarulabs/di"
)

const (
	ConfigServiceName           = "database:config"
	ConnectionServiceName       = "database:connection"
	MigrateCommandServiceName   = "database:command.migrate"
	PriceRepositoryServiceName  = "database:repository.price"
	PriceManagerServiceName     = "database:manager.price"
	ReportRepositoryServiceName = "database:repository.report"

	HostDefault    = "localhost"
	PortDefault    = "5432"
	SslModeDefault = "disable"

	hostEnv     = "APP_DB_HOST"
	portEnv     = "APP_DB_PORT"
	sslModeEnv  = "APP_DB_SSL_MODE"
	databaseEnv = "APP_DB_NAME"
	userEnv     = "APP_DB_USER"
	passwordEnv = "APP_DB_PASSWORD"
)

var Services = []di.Def{
	{
		Name:  ConfigServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cnf := &Config{
				Host:     config.OptionalEnv(hostEnv, HostDefault),
				Port:     config.OptionalEnv(portEnv, PortDefault),
				SslMode:  config.OptionalEnv(sslModeEnv, SslModeDefault),
				Database: config.RequiredEnv(databaseEnv),
				User:     config.RequiredEnv(userEnv),
				Password: config.RequiredEnv(passwordEnv),
				Timezone: ctn.Get(general.ConfigServiceName).(*general.Config).Timezone,
			}

			return cnf, nil
		},
	},
	{
		Name:  ConnectionServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			ctx := ctn.Get(general.ContextServiceName).(context.Context)
			config := ctn.Get(ConfigServiceName).(*Config)

			return db.New(ctx, config.DataSourceName())
		},
	},
	{
		Name:  MigrateCommandServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			config := ctn.Get(ConfigServiceName).(*Config)

			return migrate.New(config.DataSourceName()), nil
		},
	},
	{
		Name:  PriceRepositoryServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			conn := ctn.Get(ConnectionServiceName).(*db.Database)
			repository := price.NewRepository(conn)

			return repository, nil
		},
	},
	{
		Name:  ReportRepositoryServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			conn := ctn.Get(ConnectionServiceName).(*db.Database)
			repository := report.NewRepository(conn)

			return repository, nil
		},
	},
	{
		Name:  PriceManagerServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			manager := price.NewManager()
			manager.AddProcessor(price.NewValidateProcessor())

			return manager, nil
		},
	},
}

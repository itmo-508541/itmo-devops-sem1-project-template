package services

import (
	"context"
	"project_sem/internal/app/price"
	"project_sem/internal/app/report"
	"project_sem/internal/app/settings"
	"project_sem/internal/config"
	"project_sem/internal/database"

	"github.com/go-playground/validator/v10"
	"github.com/sarulabs/di"
)

const (
	DatabaseSettingsServiceName = "database:settings"
	ConnectionServiceName       = "database:connection"
	PriceRepositoryServiceName  = "database:repository.price"
	ReportRepositoryServiceName = "database:repository.report"

	DatabaseHostDefault    = "localhost"
	DatabasePortDefault    = "5432"
	DatabaseSslModeDefault = "disable"

	databaseHostEnv     = "APP_DB_HOST"
	databasePortEnv     = "APP_DB_PORT"
	databaseSslModeEnv  = "APP_DB_SSL_MODE"
	databaseNameEnv     = "APP_DB_NAME"
	databaseUserEnv     = "APP_DB_USER"
	databasePasswordEnv = "APP_DB_PASSWORD"
)

var DatabaseServices = []di.Def{
	{
		Name:  DatabaseSettingsServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			var generalSettings *settings.GeneralSettings

			if err := ctn.Fill(GeneralSettingsServiceName, &generalSettings); err != nil {
				return nil, err
			}
			cnf := &settings.DatabaseSettings{
				Host:     config.OptionalEnv(databaseHostEnv, DatabaseHostDefault),
				Port:     config.OptionalEnv(databasePortEnv, DatabasePortDefault),
				SslMode:  config.OptionalEnv(databaseSslModeEnv, DatabaseSslModeDefault),
				Database: config.RequiredEnv(databaseNameEnv),
				User:     config.RequiredEnv(databaseUserEnv),
				Password: config.RequiredEnv(databasePasswordEnv),
				Timezone: generalSettings.Timezone,
			}

			return cnf, nil
		},
	},
	{
		Name:  ConnectionServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			var ctx context.Context
			var config *settings.DatabaseSettings

			if err := ctn.Fill(RootContextServiceName, &ctx); err != nil {
				return nil, err
			}
			if err := ctn.Fill(DatabaseSettingsServiceName, &config); err != nil {
				return nil, err
			}

			return database.New(ctx, config.DataSourceName())
		},
	},
	{
		Name:  PriceRepositoryServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			var conn *database.Database
			var v *validator.Validate

			if err := ctn.Fill(ConnectionServiceName, &conn); err != nil {
				return nil, err
			}
			if err := ctn.Fill(ValidatorServiceName, &v); err != nil {
				return nil, err
			}

			return price.NewRepository(conn, v), nil
		},
	},
	{
		Name:  ReportRepositoryServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			var conn *database.Database

			if err := ctn.Fill(ConnectionServiceName, &conn); err != nil {
				return nil, err
			}

			return report.NewRepository(conn), nil
		},
	},
}

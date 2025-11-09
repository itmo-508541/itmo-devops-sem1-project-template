package database

import (
	"context"
	"fmt"
	"log"
	"project_sem/internal/config"
	"project_sem/internal/services/general"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/sarulabs/di"
	"github.com/spf13/cobra"
)

const (
	ConfigServiceName         = "database:config"
	ConnectionServiceName     = "database:connection"
	MigrateCommandServiceName = "database:command.migrate"

	HostDefault    = "localhost"
	PortDefault    = "5432"
	SslModeDefault = "disable"

	hostEnv     = "APP_DB_HOST"
	portEnv     = "APP_DB_PORT"
	sslModeEnv  = "APP_DB_SSL_MODE"
	databaseEnv = "APP_DB_NAME"
	userEnv     = "APP_DB_USER"
	passwordEnv = "APP_DB_PASSWORD"

	migrateUse = "migrate"
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

			db, err := pgx.Connect(ctx, config.DataSourceName())
			if err != nil {
				return nil, fmt.Errorf("pgx.Connect: %w", err)
			}

			return &Connection{db}, nil
		},
	},
	{
		Name:  MigrateCommandServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dsn := ctn.Get(ConfigServiceName).(*Config).DataSourceName()
			cmd := &cobra.Command{
				Use:   migrateUse,
				Short: "Migrate database schema",
				RunE: func(cmd *cobra.Command, args []string) error {
					m, err := migrate.New("file://migrations", dsn)
					if err != nil {
						return err
					}
					if err := m.Up(); err != nil && err != migrate.ErrNoChange {
						return err
					}
					version, _, _ := m.Version()
					log.Println("Migrated to version:", version)

					return nil
				},
			}

			return cmd, nil
		},
	},
}

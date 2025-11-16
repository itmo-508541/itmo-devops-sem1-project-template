package services

import (
	"context"
	"net/http"
	"project_sem/internal/app/command"
	"project_sem/internal/app/price"
	"project_sem/internal/app/settings"

	"github.com/sarulabs/di"
	"github.com/spf13/cobra"
)

const (
	MigrateCommandServiceName     = "command:migrate"
	StartServerCommandServiceName = "command:start-server"
	RootCommandServiceName        = "command:root"
)

var CommandServices = []di.Def{
	{
		Name:  MigrateCommandServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			var config *settings.DatabaseSettings

			if err := ctn.Fill(DatabaseSettingsServiceName, &config); err != nil {
				return nil, err
			}

			return command.NewMigrate(config.DataSourceName()), nil
		},
	},
	{
		Name:  StartServerCommandServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			var srv *http.Server
			var ctx context.Context
			var repo *price.Repository

			if err := ctn.Fill(WebServerServiceName, &srv); err != nil {
				return nil, err
			}
			if err := ctn.Fill(RootContextServiceName, &ctx); err != nil {
				return nil, err
			}
			if err := ctn.Fill(PriceRepositoryServiceName, &repo); err != nil {
				return nil, err
			}

			return command.NewStartServer(ctx, srv, repo), nil
		},
	},
	{
		Name:  RootCommandServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			var migrateCmd, startServerCmd *cobra.Command

			if err := ctn.Fill(MigrateCommandServiceName, &migrateCmd); err != nil {
				return nil, err
			}
			if err := ctn.Fill(StartServerCommandServiceName, &startServerCmd); err != nil {
				return nil, err
			}

			rootCmd := &cobra.Command{
				Short: "Final project 1st semester (Andrey Mindubaev, id:508541)",
			}
			rootCmd.AddCommand(migrateCmd)
			rootCmd.AddCommand(startServerCmd)

			return rootCmd, nil
		},
	},
}

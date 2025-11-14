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
		Build: func(ctn di.Container) (interface{}, error) {
			config := ctn.Get(DatabaseSettingsServiceName).(*settings.DatabaseSettings)

			return command.NewMigrate(config.DataSourceName()), nil
		},
	},
	{
		Name:  StartServerCommandServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			srv := ctn.Get(WebServerServiceName).(*http.Server)
			ctx := ctn.Get(RootContextServiceName).(context.Context)
			repo := ctn.Get(PriceRepositoryServiceName).(*price.Repository)

			return command.NewStartServer(ctx, srv, repo), nil
		},
	},
	{
		Name:  RootCommandServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			rootCmd := &cobra.Command{
				Short: "Final project 1st semester (Andrey Mindubaev, id:508541)",
			}
			rootCmd.AddCommand(ctn.Get(MigrateCommandServiceName).(*cobra.Command))
			rootCmd.AddCommand(ctn.Get(StartServerCommandServiceName).(*cobra.Command))

			return rootCmd, nil
		},
	},
}

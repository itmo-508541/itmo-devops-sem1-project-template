package services

import (
	"context"
	"os/signal"
	"project_sem/internal/app/settings"
	"project_sem/internal/app/validators"
	"project_sem/internal/config"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/sarulabs/di"
)

const (
	GeneralSettingsServiceName = "general:settings"
	RootContextServiceName     = "general:context"
	ValidatorServiceName       = "general:validator"

	TimezoneDefault = "Europe/Moscow"

	timezoneEnv = "APP_TIMEZONE"
)

var GeneralServices = []di.Def{
	{
		Name:  GeneralSettingsServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			cfg := &settings.GeneralSettings{
				Timezone: config.OptionalEnv(timezoneEnv, TimezoneDefault),
			}

			return cfg, nil
		},
	},
	func() di.Def {
		rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

		return di.Def{
			Name:  RootContextServiceName,
			Scope: di.App,
			Build: func(ctn di.Container) (any, error) {
				return rootCtx, nil
			},
			Close: func(obj any) error {
				stop()
				return nil
			},
		}
	}(),
	{
		Name:  ValidatorServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			v := validator.New()
			v.RegisterValidation("date", validators.DateValidator())
			v.RegisterValidation("notblank", validators.NotBlankValidator())

			return v, nil
		},
	},
}

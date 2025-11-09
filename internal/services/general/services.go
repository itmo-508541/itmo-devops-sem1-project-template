package general

import (
	"context"
	"project_sem/internal/config"

	"github.com/sarulabs/di"
)

const (
	ConfigServiceName  = "general:config"
	ContextServiceName = "general:context"

	TimezoneDefault = "Europe/Moscow"

	timezoneEnv = "APP_TIMEZONE"
)

var Services = []di.Def{
	{
		Name:  ConfigServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := &Config{
				Timezone: config.OptionalEnv(timezoneEnv, TimezoneDefault),
			}

			return cfg, nil
		},
	},
	func() di.Def {
		context, cancel := context.WithCancel(context.Background())

		return di.Def{
			Name:  ContextServiceName,
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				return context, nil
			},
			Close: func(obj interface{}) error {
				cancel()
				return nil
			},
		}
	}(),
}

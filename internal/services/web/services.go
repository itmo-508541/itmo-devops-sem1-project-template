package web

import (
	"context"
	_ "embed"
	"net/http"
	"project_sem/internal/app/assets"
	"project_sem/internal/app/handlers/load"
	"project_sem/internal/app/handlers/save"
	"project_sem/internal/config"
	"project_sem/internal/models/price"
	"project_sem/internal/server"
	"project_sem/internal/server/command/start"
	"project_sem/internal/services/database"
	"project_sem/internal/services/general"

	"github.com/sarulabs/di"
)

const (
	ConfigServiceName             = "web:config"
	LoadHandlerServiceName        = "web:handler.load"
	SaveHandlerServiceName        = "web:handler.save"
	ServeMuxServiceName           = "web:router"
	ServerServiceName             = "web:server"
	StartServerCommandServiceName = "web:start-server"

	HostDefault = "0.0.0.0"
	PortDefault = "8080"

	hostEnv = "APP_HOST"
	portEnv = "APP_PORT"
)

var Services = []di.Def{
	{
		Name:  ConfigServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := &Config{
				Host: config.OptionalEnv(hostEnv, HostDefault),
				Port: config.OptionalEnv(portEnv, PortDefault),
			}

			return cfg, nil
		},
	},
	{
		Name:  LoadHandlerServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			handler := load.New()

			return handler, nil
		},
	},
	{
		Name:  SaveHandlerServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			manager := ctn.Get(database.PriceManagerServiceName).(*price.Manager)
			handler := save.New(manager)

			return handler, nil
		},
	},
	{
		Name:  ServeMuxServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			loadHandler := ctn.Get(LoadHandlerServiceName).(*load.Handler)
			saveHandler := ctn.Get(SaveHandlerServiceName).(*save.Handler)

			mux := http.NewServeMux()
			mux.Handle("GET /api/v0/prices", loadHandler)
			mux.Handle("POST /api/v0/prices", saveHandler)

			mux.Handle("/panic", server.PanicRecoveryMiddleware(
				func(http.ResponseWriter, *http.Request) {
					panic("panic to recover")
				},
			))
			mux.Handle("/favicon.ico", http.FileServer(http.FS(assets.FaviconFS)))
			mux.Handle("/", http.FileServer(http.FS(assets.IndexFS)))

			return mux, nil
		},
	},
	{
		Name:  ServerServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			mux := ctn.Get(ServeMuxServiceName).(*http.ServeMux)
			config := ctn.Get(ConfigServiceName).(*Config)

			return server.New(mux, config.Addr()), nil
		},
	},
	{
		Name:  StartServerCommandServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			srv := ctn.Get(ServerServiceName).(*http.Server)
			ctx := ctn.Get(general.ContextServiceName).(context.Context)
			repo := ctn.Get(database.PriceRepositoryServiceName).(*price.Repository)

			return start.New(ctx, srv, repo), nil
		},
	},
}

package services

import (
	"context"
	_ "embed"
	"net/http"
	"project_sem/internal/app/assets"
	"project_sem/internal/app/command"
	"project_sem/internal/app/price"
	"project_sem/internal/app/report"
	"project_sem/internal/app/server"
	"project_sem/internal/app/settings"
	"project_sem/internal/config"

	"github.com/sarulabs/di"
)

const (
	WebSettingsServiceName        = "web:settings"
	LoadHandlerServiceName        = "web:handler.load"
	SaveHandlerServiceName        = "web:handler.save"
	ServeMuxServiceName           = "web:router"
	ServerServiceName             = "web:server"
	StartServerCommandServiceName = "web:start-server"

	WebHostDefault = "0.0.0.0"
	WebPortDefault = "8080"

	webHostEnv = "APP_HOST"
	webPortEnv = "APP_PORT"
)

var WebServices = []di.Def{
	{
		Name:  WebSettingsServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := &settings.WebSettings{
				Host: config.OptionalEnv(webHostEnv, WebHostDefault),
				Port: config.OptionalEnv(webPortEnv, WebPortDefault),
			}

			return cfg, nil
		},
	},
	{
		Name:  LoadHandlerServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			reportR := ctn.Get(ReportRepositoryServiceName).(*report.Repository)
			handler := server.NewLoadHandler(reportR)

			return handler, nil
		},
	},
	{
		Name:  SaveHandlerServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			manager := ctn.Get(PriceManagerServiceName).(*price.Manager)
			priceR := ctn.Get(PriceRepositoryServiceName).(*price.Repository)
			reportR := ctn.Get(ReportRepositoryServiceName).(*report.Repository)

			handler := server.NewSaveHandler(manager, priceR, reportR)

			return handler, nil
		},
	},
	{
		Name:  ServeMuxServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			loadHandler := ctn.Get(LoadHandlerServiceName).(http.HandlerFunc)
			saveHandler := ctn.Get(SaveHandlerServiceName).(http.HandlerFunc)

			mux := http.NewServeMux()
			mux.Handle("GET /api/v0/prices", server.PanicRecoveryMiddleware(loadHandler))
			mux.Handle("POST /api/v0/prices", server.PanicRecoveryMiddleware(saveHandler))
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
			config := ctn.Get(WebSettingsServiceName).(*settings.WebSettings)

			return server.NewWebServer(mux, config), nil
		},
	},
	{
		Name:  StartServerCommandServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			srv := ctn.Get(ServerServiceName).(*http.Server)
			ctx := ctn.Get(RootContextServiceName).(context.Context)
			repo := ctn.Get(PriceRepositoryServiceName).(*price.Repository)

			return command.NewStartServer(ctx, srv, repo), nil
		},
	},
}

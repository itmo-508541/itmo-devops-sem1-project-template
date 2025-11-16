package services

import (
	"net/http"
	"project_sem/internal/app/assets"
	"project_sem/internal/app/price"
	"project_sem/internal/app/report"
	"project_sem/internal/app/server"
	"project_sem/internal/app/settings"
	"project_sem/internal/config"

	"github.com/go-playground/validator/v10"
	"github.com/sarulabs/di"
)

const (
	WebSettingsServiceName = "web:settings"
	LoadHandlerServiceName = "web:handler.load"
	SaveHandlerServiceName = "web:handler.save"
	ServeMuxServiceName    = "web:router"
	WebServerServiceName   = "web:server"

	WebHostDefault = "0.0.0.0"
	WebPortDefault = "8080"

	webHostEnv = "APP_HOST"
	webPortEnv = "APP_PORT"
)

var WebServices = []di.Def{
	{
		Name:  WebSettingsServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
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
		Build: func(ctn di.Container) (any, error) {
			var v *validator.Validate
			var reportR *report.Repository

			if err := ctn.Fill(ValidatorServiceName, &v); err != nil {
				return nil, err
			}
			if err := ctn.Fill(ReportRepositoryServiceName, &reportR); err != nil {
				return nil, err
			}

			handler := server.NewLoadHandler(reportR, v)

			return handler, nil
		},
	},
	{
		Name:  SaveHandlerServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			var priceR *price.Repository
			var reportR *report.Repository

			if err := ctn.Fill(PriceRepositoryServiceName, &priceR); err != nil {
				return nil, err
			}
			if err := ctn.Fill(ReportRepositoryServiceName, &reportR); err != nil {
				return nil, err
			}

			handler := server.NewSaveHandler(priceR, reportR)

			return handler, nil
		},
	},
	{
		Name:  ServeMuxServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			var loadHandler, saveHandler http.HandlerFunc

			if err := ctn.Fill(LoadHandlerServiceName, &loadHandler); err != nil {
				return nil, err
			}
			if err := ctn.Fill(SaveHandlerServiceName, &saveHandler); err != nil {
				return nil, err
			}

			mux := http.NewServeMux()
			mux.Handle("GET /api/v0/prices", server.PanicRecoveryMiddleware(loadHandler))
			mux.Handle("POST /api/v0/prices", server.PanicRecoveryMiddleware(saveHandler))
			mux.Handle("/favicon.ico", http.FileServer(http.FS(assets.FaviconFS)))
			mux.Handle("/", http.FileServer(http.FS(assets.IndexFS)))

			return mux, nil
		},
	},
	{
		Name:  WebServerServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			var mux *http.ServeMux
			var config *settings.WebSettings

			if err := ctn.Fill(ServeMuxServiceName, &mux); err != nil {
				return nil, err
			}
			if err := ctn.Fill(WebSettingsServiceName, &config); err != nil {
				return nil, err
			}

			return server.NewWebServer(mux, config), nil
		},
	},
}

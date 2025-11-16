package server

import (
	"net/http"
	"time"

	"project_sem/internal/app/settings"
)

func NewWebServer(mux *http.ServeMux, cfg *settings.WebSettings) *http.Server {
	return &http.Server{
		Handler:      mux,
		Addr:         cfg.Addr(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

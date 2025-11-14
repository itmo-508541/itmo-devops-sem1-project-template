package server

import (
	"net/http"
	"project_sem/internal/app/settings"
	"time"
)

func NewWebServer(mux *http.ServeMux, cfg *settings.WebSettings) *http.Server {
	return &http.Server{
		Handler:      mux,
		Addr:         cfg.Addr(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

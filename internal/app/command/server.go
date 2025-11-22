package command

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"project_sem/internal/app/assets"
	"project_sem/internal/app/server"
	"project_sem/internal/app/settings"
	"project_sem/internal/database"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

const startServerUse = "start-server"

func NewStartServer() *cobra.Command {
	return &cobra.Command{
		Use:   startServerUse,
		Short: "Start web-server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			conn, err := database.New(ctx, settings.DatabaseSourceName())
			if err != nil {
				return err
			}

			loadHandler := server.NewLoadHandler(conn)
			saveHandler := server.NewSaveHandler(conn)

			mux := http.NewServeMux()
			mux.Handle("GET /api/v0/prices", server.PanicRecoveryMiddleware(loadHandler))
			mux.Handle("POST /api/v0/prices", server.PanicRecoveryMiddleware(saveHandler))
			mux.Handle("GET /favicon.ico", http.FileServer(http.FS(assets.FaviconFS)))
			// мне нужна web-страница для ручного тестирования
			mux.Handle("GET /", http.FileServer(http.FS(assets.IndexFS)))

			srv := &http.Server{
				Handler:      mux,
				Addr:         settings.WebServerAddr(),
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			log.Println("Starting web-server...")
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Println(fmt.Errorf("srv.ListenAndServe: %w", err))
				}
			}()
			log.Printf("Listening on %s", srv.Addr)
			<-ctx.Done()

			log.Println("Stopping Web-server...")
			err = srv.Shutdown(context.Background()) //nolint:contextcheck
			if err != nil {
				log.Println(fmt.Errorf("srv.Shutdown: %w", err))
			}
			log.Println("Web-server stopped")

			return err
		},
	}
}

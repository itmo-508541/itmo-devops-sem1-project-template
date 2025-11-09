package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project_sem/internal/config"
	"syscall"
	"time"

	"github.com/sarulabs/di"
	"github.com/spf13/cobra"
)

const (
	ConfigServiceName             = "web:config"
	ServerServiceName             = "web:server"
	StartServerCommandServiceName = "web:start-server"

	PortDefault = "8080"

	portEnv = "APP_PORT"

	startServerUse = "start-server"
)

var Services = []di.Def{
	{
		Name:  ConfigServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := &Config{
				Port: config.OptionalEnv(portEnv, PortDefault),
			}

			return cfg, nil
		},
	},
	{
		Name:  ServerServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			mux := http.NewServeMux()
			// @todo добавить новые handlers

			port := ctn.Get(ConfigServiceName).(*Config).Port
			// @todo тут куда-то нужно добавить контекст, наверное?? хз
			srv := &http.Server{
				Handler:      mux,
				Addr:         fmt.Sprintf("0.0.0.0:%s", port),
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			return srv, nil
		},
	},
	{
		Name:  StartServerCommandServiceName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cmd := &cobra.Command{
				Use:   startServerUse,
				Short: "Start web-server",
				// @todo наверное лучше это куда-то перенести =)
				RunE: func(cmd *cobra.Command, args []string) error {
					srv := ctn.Get(ServerServiceName).(*http.Server)

					// @see https://github.com/sarulabs/di-example/blob/master/main.go
					log.Printf("Listening on %s", srv.Addr)
					go func() {
						if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
							log.Println(err.Error())
						}
					}()

					stop := make(chan os.Signal, 1)
					signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
					<-stop

					ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
					defer cancel()

					log.Println("Stopping the http server")

					return srv.Shutdown(ctx)
				},
			}

			return cmd, nil
		},
	},
}

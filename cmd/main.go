package main

import (
	"fmt"
	"log"
	"project_sem/internal/services/application"
	"project_sem/internal/services/database"
	"project_sem/internal/services/general"
	"project_sem/internal/services/web"

	"github.com/sarulabs/di"
	"github.com/spf13/cobra"
)

func main() {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Fatal(panicErr)
		}
	}()

	builder, err := di.NewBuilder(di.App)
	if err != nil {
		log.Fatal(fmt.Errorf("di.NewBuilder: %w", err))
	}
	if err := builder.Add(general.Services...); err != nil {
		log.Fatal(fmt.Errorf("builder.Add(general.Services): %w", err))
	}
	if err := builder.Add(application.Services...); err != nil {
		log.Fatal(fmt.Errorf("builder.Add(application.Services): %w", err))
	}
	if err := builder.Add(database.Services...); err != nil {
		log.Fatal(fmt.Errorf("builder.Add(database.Services): %w", err))
	}
	if err := builder.Add(web.Services...); err != nil {
		log.Fatal(fmt.Errorf("builder.Add(web.Services): %w", err))
	}
	ctn := builder.Build()

	rootCmd := ctn.Get(application.CommandRootServiceName).(*cobra.Command)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(fmt.Errorf("rootCmd.Execute: %w", err))
	}
}

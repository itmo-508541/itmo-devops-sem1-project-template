package main

import (
	"fmt"
	"log"

	"github.com/sarulabs/di"
	"github.com/spf13/cobra"

	"project_sem/internal/app/services"
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
	for _, services := range [][]di.Def{services.GeneralServices, services.CommandServices, services.DatabaseServices, services.WebServices} {
		if err := builder.Add(services...); err != nil {
			log.Fatal(fmt.Errorf("builder.Add: %w", err))
		}
	}

	ctn := builder.Build()
	defer ctn.DeleteWithSubContainers()

	rootCmd := ctn.Get(services.RootCommandServiceName).(*cobra.Command)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(fmt.Errorf("rootCmd.Execute: %w", err))
	}
}

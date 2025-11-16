package main

import (
	"fmt"
	"log"
	"project_sem/internal/app/services"

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
		panic(fmt.Errorf("di.NewBuilder: %w", err))
	}
	for _, services := range [][]di.Def{services.GeneralServices, services.CommandServices, services.DatabaseServices, services.WebServices} {
		if err := builder.Add(services...); err != nil {
			panic(fmt.Errorf("builder.Add: %w", err))
		}
	}

	ctn := builder.Build()
	defer func() {
		err := ctn.DeleteWithSubContainers()
		if err != nil {
			log.Println(fmt.Errorf("ctn.DeleteWithSubContainers: %w", err))
		}
	}()

	var rootCmd *cobra.Command
	if err := ctn.Fill(services.RootCommandServiceName, &rootCmd); err != nil {
		panic(err)
	}
	if err := rootCmd.Execute(); err != nil {
		log.Println(fmt.Errorf("rootCmd.Execute: %w", err))
	}
}

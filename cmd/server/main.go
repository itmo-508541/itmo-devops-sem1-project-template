package main

import (
	"log"
	"project_sem/internal/app"

	_ "github.com/gorilla/mux"
)

func main() {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Fatal(panicErr)
		}
	}()
	if err := app.Init(); err != nil {
		panic(err)
	}
}

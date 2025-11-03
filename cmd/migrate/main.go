package main

import (
	"fmt"
	"log"

	"project_sem/internal/env"
	"project_sem/internal/infrastructure/database"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Fatal(panicErr)
		}
	}()

	if err := env.Load(); err != nil {
		panic(err)
	}

	m, err := migrate.New("file://migrations", database.DataSourceName())
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		fmt.Println(err)
	} else {
		version, _, _ := m.Version()
		fmt.Println("Migrated to", version)
	}
}

package app

import (
	"project_sem/internal/database"
	"project_sem/internal/env"
)

var connection *database.Connection

func Init() error {
	var err error

	if err = env.Init(); err != nil {
		return err
	}

	connection, err = database.New(env.DataSourceName())
	if err != nil {
		return err
	}

	return nil
}

func DatabaseConnection() *database.Connection {
	return connection
}

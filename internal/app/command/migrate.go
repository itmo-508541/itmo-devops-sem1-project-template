package command

import (
	"log"
	"project_sem/internal/app/migrations"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/cobra"
)

const commandUse = "migrate"

func NewMigrate(dsn string) *cobra.Command {
	return &cobra.Command{
		Use:   commandUse,
		Short: "Migrate database schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			schema, err := iofs.New(migrations.Schema, migrations.SchemaPath)
			if err != nil {
				return err
			}
			m, err := migrate.NewWithSourceInstance("iofs", schema, dsn)
			if err != nil {
				return err
			}
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				return err
			}
			version, _, _ := m.Version()
			log.Printf("Migrated to version #%d\n", version)

			return nil
		},
	}
}

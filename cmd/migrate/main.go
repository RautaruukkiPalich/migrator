package main

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rautaruukkipalich/migrator/pkg/confighelper"
)

func main() {
	cfg := confighelper.MustLoadConfig()

	m, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.MigrationPath),
		cfg.Database.DBUri,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
}

package main

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rautaruukkipalich/migrator/pkg/confighelper"
	"github.com/rautaruukkipalich/migrator/pkg/dbhelper"
)

const (
	defaultTestCountRows = 25000
)

func main() {
	cfg := confighelper.MustLoadConfig()

	dbURI, driver, err := dbhelper.GetURIAndDriverFromCfg(&cfg.Database)
	if err != nil {
		panic(err)
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.MigrationPath),
		dbURI,
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

	if cfg.TestCountRows == 0 {
		cfg.TestCountRows = defaultTestCountRows
	}

	err = FillDB(driver, dbURI, cfg.TestCountRows)
	if err != nil {
		panic(err)
	}

	fmt.Println("migrations apply successfully")
}

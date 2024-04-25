package migrator

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func newDatabase(cfg *DBConfig) (*sql.DB, error) {
	op := "init database"
	
	var path string
	var driver string

	switch cfg.Driver {
	case Postgres:
		path = fmt.Sprintf(
			"%s://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.Driver,
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
		)
		driver = "postgres"
	// case MySQL:
	// 	driver = "mysql"
	// case SQLite:
	// 	driver = "sqlite"
	default:

	}
	if driver == "" {
		return nil, fmt.Errorf("%s is not a valid driver", cfg.Driver)
	}

	db, err := sql.Open(driver, path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

package dbhelper

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rautaruukkipalich/migrator/config"
)

func GetDBConnection(cfg *config.DatabaseConfig) (*sqlx.DB, error) {

	if cfg.Driver == "" {
		return nil, ErrInvalidDriver
	}

	if cfg.DBUri == "" {
		return nil, ErrInvalidDatabaseURI
	}

	db, err := sqlx.Open(string(cfg.Driver), cfg.DBUri)
	if err != nil {
		return nil, ErrFailedToCreateConnection
	}

	if err := db.Ping(); err != nil {
		return nil, ErrFailedToConnectDB
	}

	return db, nil
}
package dbhelper

import (
	"errors"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/rautaruukkipalich/migrator/config"
)

func GetURIAndDriverFromCfg(cfg *config.DatabaseConfig) (string, string, error) {
	var dbURI string

	switch cfg.Driver {
	case config.Postgres:
		dbURI = fmt.Sprintf(
			"%s://%s:%s@%s:%v/%s?sslmode=disable",
			cfg.Driver,	cfg.User, cfg.Password,	cfg.Host, cfg.Port, cfg.DBName,
		)
	default:
		// TODO add sqlite3 support
		return "", "", errors.New("use postgres")
	}

	return dbURI, string(cfg.Driver), nil
}

func GetDBConnection(dbURI, driver string) (*sqlx.DB, error) {

	if driver == "" {
		return nil, ErrInvalidDriver
	}

	db, err := sqlx.Open(driver, dbURI)
	if err != nil {
		return nil, ErrFailedToCreateConnection
	}

	if err := db.Ping(); err != nil {
		return nil, ErrFailedToConnectDB
	}

	return db, nil
}
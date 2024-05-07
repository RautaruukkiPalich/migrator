package migrator

import (
	"fmt"
	"math"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

const (
	batchSize = 10000
)

func newDatabase(cfg *DBConfig) (*sqlx.DB, error) {

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
		return nil, ErrInvalidDriver
	}

	db, err := sqlx.Open(driver, path)
	if err != nil {
		return nil, ErrFailedToCreateConnection
	}

	if err := db.Ping(); err != nil {
		return nil, ErrFailedToConnectDB
	}

	return db, nil
}

func (m *migrator) MigrateFromDB(table string) error {
	tx := m.donor.MustBegin()
	defer tx.Rollback()

	var count []any
	err := tx.Select(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s", table))
	if err != nil {
		return err
	}

	times := int(math.Ceil(float64(count[0].(int64)) / batchSize))

	for i := 0; i < times; i++ {
		limit := batchSize
		offset := i*batchSize

		rows, err := tx.Queryx(
			fmt.Sprintf(`SELECT * FROM %s LIMIT $1 OFFSET $2`, table),
			limit,
			offset,
		)
		if err != nil {
			return fmt.Errorf("select from db err: %w", err)
		}
		
		err = m.SendMessages(table, rows)

		if err != nil {
			return fmt.Errorf("err send to kafka: %w", err)
		}	
	}

	tx.Commit()
	return nil
}

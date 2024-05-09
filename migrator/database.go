package migrator

import (
	"fmt"
	"math"

	"github.com/jmoiron/sqlx"
	"github.com/rautaruukkipalich/migrator/config"
	"github.com/rautaruukkipalich/migrator/pkg/dbhelper"
)

func newDatabase(cfg *config.DatabaseConfig) (*sqlx.DB, error) {

	dbURI, driver, err := dbhelper.GetURIAndDriverFromCfg(cfg)
	if err != nil {
		return nil, err
	}
	
	return dbhelper.GetDBConnection(dbURI, driver)
}

func (m *migrator) MigrateFromDB(table string) error {
	tx := m.database.MustBegin()
	defer tx.Rollback()

	rowCount, err := m.getRowsCount(table, tx) 
	if err != nil {
		return err
	}

	iterations := m.getIterationsRange(rowCount)

	stmt, err := tx.Preparex(
		fmt.Sprintf(`SELECT * FROM %s LIMIT $1 OFFSET $2`, table),
	)
	if err != nil {
		return err
	}

	for i := 0; i < iterations; i++ {
		rows, err := m.getRows(table, stmt, i)
		if err != nil {
			return err
		}

		if err = m.SendMessages(table, rows); err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}
	
func (m *migrator) getRowsCount(table string, tx *sqlx.Tx) (int, error) {

	var count []any
	err := tx.Select(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s", table))
	if err != nil {
		return 0, err
	}
	return int(count[0].(int64)), nil
}


func (m *migrator) getRows(table string, stmt *sqlx.Stmt, iter int) (*sqlx.Rows, error) {
	// tablename validate in m.Migrate func
	return stmt.Queryx(m.batchSize, iter * m.batchSize)
}

func (m *migrator) getIterationsRange(count int) int {	
	return int(
		math.Ceil(
			float64(count) / float64(m.batchSize),
		),
	)
}

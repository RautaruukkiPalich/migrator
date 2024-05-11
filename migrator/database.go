package migrator

import (
	_ "embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rautaruukkipalich/migrator/config"
	"github.com/rautaruukkipalich/migrator/pkg/dbhelper"
)

var (
	//go:embed queries/select_rows.sql
	selectRows string
)


func newDatabase(cfg *config.DatabaseConfig) (*sqlx.DB, error) {

	dbURI, driver, err := dbhelper.GetURIAndDriverFromCfg(cfg)
	if err != nil {
		return nil, err
	}
	
	return dbhelper.GetDBConnection(dbURI, driver)
}

func (m *migrator) MigrateTable(table string) error {
	tx := m.database.MustBegin()
	//nolint:all
	defer tx.Rollback()

	stmt, err := tx.Preparex(fmt.Sprintf(selectRows, table))
	if err != nil {
		return err
	}

	rows, err := stmt.Queryx()
	if err != nil {
		return err
	}

	if err = m.SendToBroker(table, rows); err != nil {
		return err
	}

	//nolint:all
	tx.Commit()
	return nil
}
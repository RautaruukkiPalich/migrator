package sql

import (
	"context"
	_ "embed"

	"github.com/jmoiron/sqlx"
	"github.com/rautaruukkipalich/migrator/config"
	"github.com/rautaruukkipalich/migrator/migrator/database"
	"github.com/rautaruukkipalich/migrator/pkg/dbhelper"
)

type Database struct {
	conn *sqlx.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	conn, err := dbhelper.GetDBConnection(&cfg.Database)
	if err != nil {
		return nil, err
	}
	return &Database{conn: conn}, nil
}

func (db *Database) Close() {
	db.conn.Close()
}

func (db *Database) GetRows(ctx context.Context, selectQuery string) (chan database.Row, error) {

	rows, err := db.conn.Queryx(selectQuery)
	if err != nil {
		return nil, database.ErrGetRows
	}

	ch := make(chan database.Row)

	go func(rows *sqlx.Rows){
		defer rows.Close()
		defer close(ch)

		for rows.Next() {
			rowMap := make(map[string]any)
			err := rows.MapScan(rowMap)
	
			select {
			case ch <- database.Row{
				Row: rowMap,
				Err: database.ErrParseRow,
			}: if err != nil {
				return
			}
			case <-ctx.Done():
				return
			}
		}
	}(rows)

	return ch, nil
}

package sql

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

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

	selectQuery = db.TrimSpaces(selectQuery)

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
			if err != nil {
				err = database.ErrGetRows
			}
	
			select {
			case ch <- database.Row{
				Row: rowMap,
				Err: err,
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

func (db *Database) TrimSpaces(query string) string {
	return strings.TrimSpace(query)
}

//если буду делить на несколько запросов чтоб избежать проблем со слишком большим количеством памяти для множества sql.Rows
func (db *Database) GetCountRows(query string) int32 {
	query = strings.TrimRight(query, ";")
	query = fmt.Sprintf("select count(*) from (%s) as subquery;", query)

	res := db.conn.QueryRowx(query)
	dict := make(map[string]any)
	res.MapScan(dict)

	return int32(dict["count"].(int64))
}

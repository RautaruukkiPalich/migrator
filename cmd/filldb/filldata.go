package main

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/rautaruukkipalich/migrator/config"
	"github.com/rautaruukkipalich/migrator/pkg/dbhelper"
)

var (
	//go:embed queries/query.sql
	insertQuery string
)

func fillDBTestData(cfg *config.Config) error {

	db, err := dbhelper.GetDBConnection(&cfg.Database)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	tx := db.MustBegin()
	//nolint:all
	defer tx.Rollback()

	stmt, err := tx.Preparex(insertQuery)
	if err != nil {
		return err
	}

	for i := 1; i <= cfg.TestCountRows; i++ {
		username := fmt.Sprintf("user_%d", i)
		createdAt := time.Now().UTC()
		stmt.MustExec(username, createdAt)
	}

	return tx.Commit()
}

package main

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	username  string
	createdAt time.Time
}

func FillDB(driver string, path string, countRows int) error {

	db, err := sql.Open(driver, path)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	for i := 1; i <= countRows; i++ {
		user := User{
			username:  fmt.Sprintf("user_%d", i),
			createdAt: time.Now().UTC(),
		}
		saveToDB(db, user)
	}

	return nil
}

func saveToDB(db *sql.DB, user User) {
	stmt, err := db.Prepare(
		`INSERT 
		INTO donor (username, created_at)
		VALUES ($1, $2)`,
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	stmt.Exec(user.username, user.createdAt)
}

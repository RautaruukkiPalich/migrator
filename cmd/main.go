package main

import (
	"fmt"
	"github.com/rautaruukkipalich/migrator/migrator"
)

func main() {
	if err := migrator.Migrate(
		migrator.DBConfig{
			Driver: migrator.Postgres,
			Host:     "localhost",
			Port:     5441,
			Username: "postgres",
			Password: "postgres",
			DBName:   "postgres",
		},
		migrator.DBConfig{
			Driver: migrator.Postgres,
			Host:     "localhost",
			Port:     5440,
			Username: "postgres",
			Password: "postgres",
			DBName:   "postgres",
		},
		migrator.KafkaConfig{
			Host: "localhost",
			Port: 29092,
		},
	); err != nil {
		fmt.Printf("Failed to migrate: %v\n", err)
	}
	fmt.Println("success")
}

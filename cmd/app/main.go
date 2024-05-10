package main

import (
	"log"

	"github.com/rautaruukkipalich/migrator/migrator"
	"github.com/rautaruukkipalich/migrator/pkg/confighelper"
)

func main() {
	cfg := confighelper.MustLoadConfig()

	m, err := migrator.New(
		&cfg.Database,
		&cfg.Kafka, 
		cfg.BatchSize,
	)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v\n", err)
	}
	defer m.Close()

	for _, table := range cfg.Tables {
		err = m.Migrate(table)
		if err != nil {
			log.Printf("Failed to migrate table `%s`: %v\n", table, err)
		} else {
			log.Printf("Successful migrated table `%s`", table)
		}
	}

	log.Print("and its gone")
}

package main

import (
	"log"

	"github.com/rautaruukkipalich/migrator/migrator"
	"github.com/rautaruukkipalich/migrator/pkg/confighelper"
)

func main() {
	cfg := confighelper.MustLoadConfig()

	m, err := migrator.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v\n", err)
	}
	defer m.Close()

	err = m.Migrate(cfg.OutputTablename)
	if err != nil {
		log.Printf("Failed to migrate table `%s`: %v\n", cfg.OutputTablename, err)
	} else {
		log.Printf("Successful migrated table `%s`", cfg.OutputTablename)
	}

	log.Print("and its gone")
}

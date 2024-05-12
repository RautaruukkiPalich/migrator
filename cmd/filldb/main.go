package main

import (
	"fmt"

	"github.com/rautaruukkipalich/migrator/pkg/confighelper"
)

const (
	defaultTestCountRows = 25000
)

func main() {
	cfg := confighelper.MustLoadConfig()

	if cfg.TestCountRows == 0 {
		cfg.TestCountRows = defaultTestCountRows
	}

	err := fillDBTestData(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println("migrations apply successfully")
}
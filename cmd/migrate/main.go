package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilyakaznacheev/cleanenv"
)

type Driver string

const (
	countRows = 45000
	Postgres Driver = "postgres"
)

type Config struct {
	Database      DatabaseConfig `yaml:"database" env_required:"true"`
	MigrationPath string         `yaml:"migration_path" env_required:"true"`
}

type DatabaseConfig struct {
	Driver   Driver `yaml:"driver"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"db_name"`
}

func main() {
	cfg := MustLoadConfig()

	var dbURI string

	switch cfg.Database.Driver {
	case Postgres:
		dbURI = fmt.Sprintf(
			"%s://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.Database.Driver,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		)
	default:
		// TODO add sqlite3 support
		panic("use postgres")
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.MigrationPath),
		dbURI,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	err = fillDB(cfg.Database.Driver, dbURI)
	if err != nil {
		panic(err)
	}

	fmt.Println("migrations apply successfully")
}

func MustLoadConfig() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("can not parse config")
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	return res
}

type Data struct {
	username string
	created_at time.Time
}

func fillDB(driver Driver, path string) error{

	db, err := sql.Open(string(driver), path)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	for i := 1; i <= countRows; i++ {
		data := Data{
			username: fmt.Sprintf("user_%d", i),
			created_at: time.Now().UTC(),
		}
		saveToDB(db, data)
	}

	return nil
}

func saveToDB(db *sql.DB, data Data) {
	stmt, err := db.Prepare(
		`INSERT 
		INTO donor (username, created_at)
		VALUES ($1, $2)`,
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	stmt.Exec(data.username, data.created_at)
}

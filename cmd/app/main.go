package main

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rautaruukkipalich/migrator/migrator"
)

type Driver string

const (
	Postgres Driver = "postgres"
)

type Config struct {
	Database DatabaseConfig `yaml:"database" env_required:"true"`
	Kafka    KafkaConfig    `yaml:"kafka" env_required:"true"`
}

type DatabaseConfig struct {
	Driver   Driver `yaml:"driver"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"db_name"`
}

type KafkaConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func main() {
	cfg := MustLoadConfig()

	dbcfg := &migrator.DBConfig{
		Driver:   migrator.Driver(cfg.Database.Driver),
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Username: cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	}

	kafkacfg := &migrator.KafkaConfig{
		Host: cfg.Kafka.Host,
		Port: cfg.Kafka.Port,
	}

	m, err := migrator.New(dbcfg, kafkacfg)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v\n", err)
	}
	defer m.Close()

	dbs := []string{"donor", "postgres;drop table users", "postgres drop table users"}

	for _, db := range dbs {
		err = m.Migrate(db)
		if err != nil {
			log.Printf("Failed to migrate db `%s`: %v\n", db, err)
		} else {
			log.Printf("Suscessfully migrated `%s`", db)
		}
	}

	log.Print("and its gone")
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

package migrator

import (
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/jmoiron/sqlx"
)

type Driver string

const (
	Postgres Driver = "postgres"
	SQLite   Driver = "sqlite"
	MySQL    Driver = "mysql"
)

type DBConfig struct {
	Driver   Driver
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

type KafkaConfig struct {
	Host string
	Port int
}

type migrator struct {
	donor  *sqlx.DB
	broker *kafka.Writer
}

type Migrator interface {
	Migrate(table string) error
	Close()
}

func New(
	donorConfig *DBConfig,
	kafkaConfig *KafkaConfig,
) (Migrator, error) {
	donor, err := newDatabase(donorConfig)
	if err != nil {
		return nil, err
	}
	broker := newBroker(kafkaConfig)

	return &migrator{
		donor:  donor,
		broker: broker,
	}, nil
}

func (m *migrator) Migrate(table string) error {
	if err := checkTable(table); err != nil {
		return err
	}
	return m.MigrateFromDB(table)
}

func (m *migrator) Close() {
	m.broker.Close()
	m.donor.Close()
}

func checkTable(table string) error {
	if len(strings.Split(table, " ")) != 1 {
		return ErrInvalidTablename
	}
	if len(strings.Split(table, ";")) != 1 {
		return ErrInvalidTablename
	}
	if len(strings.Split(table, ",")) != 1 {
		return ErrInvalidTablename
	}
	if strings.Contains(table, "drop") {
		return ErrInvalidTablename
	}
	return nil
}

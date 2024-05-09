package migrator

import (
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/rautaruukkipalich/migrator/config"
	"github.com/segmentio/kafka-go"
)

type migrator struct {
	database  *sqlx.DB
	broker    *kafka.Writer
	batchSize int
}

type Migrator interface {
	Migrate(table string) error
	Close()
}

func New(
	cfgDB *config.DatabaseConfig,
	cfgKafka *config.KafkaConfig,
	batchSize int,
) (Migrator, error) {

	database, err := newDatabase(cfgDB)
	if err != nil {
		return nil, err
	}
	broker := newBroker(cfgKafka)

	if batchSize == 0 {
		batchSize = 10000
	}

	return &migrator{
		database:  database,
		broker:    broker,
		batchSize: batchSize,
	}, nil
}

func (m *migrator) Migrate(table string) error {
	if err := validateTable(table); err != nil {
		return err
	}
	return m.MigrateFromDB(table)
}

func (m *migrator) Close() {
	m.broker.Close()
	m.database.Close()
}

func validateTable(table string) error {
	if strings.Contains(table, "drop") ||
		strings.Contains(table, " ") ||
		strings.Contains(table, ";") ||
		strings.Contains(table, ",") ||
		strings.Contains(table, ".") {
		return ErrInvalidTablename
	}
	return nil
}

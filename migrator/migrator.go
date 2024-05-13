package migrator

import (
	"context"

	"github.com/rautaruukkipalich/migrator/config"
	kafka_broker "github.com/rautaruukkipalich/migrator/migrator/broker/kafka"
	"github.com/rautaruukkipalich/migrator/migrator/database"
	sql_database "github.com/rautaruukkipalich/migrator/migrator/database/sql"
	"github.com/rautaruukkipalich/migrator/resources"
)

type Broker interface {
	SendMessages(ctx context.Context, ch chan database.Row, tablename string) error
	Close()
}

type Database interface {
	GetRows(ctx context.Context, query string) (chan database.Row, error)
	Close()
}

type migrator struct {
	database Database
	broker   Broker
}

type Migrator interface {
	Migrate(table string) error
	Close()
}

func New(
	cfg *config.Config,
) (Migrator, error) {
	
	if cfg.BatchSize == 0 {
		cfg.BatchSize = 10000
	}

	db, err := sql_database.NewDatabase(cfg)
	if err != nil {
		return nil, err
	}
	broker, err := kafka_broker.NewBroker(cfg)
	if err != nil {
		return nil, err
	}

	return &migrator{
		database: db,
		broker:   broker,
	}, nil
}

func (m *migrator) Migrate(tablename string) error {
	if err := validateTableName(tablename); err != nil {
		return err
	}
	return m.MigrateTable(tablename)
}

func (m *migrator) Close() {
	m.broker.Close()
	m.database.Close()
}

func (m *migrator) MigrateTable(tablename string) error {

	selectQuery := resources.QuerySelect

	if err := validateSelectQuery(selectQuery); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	ch, err := m.database.GetRows(ctx, selectQuery)
	if err != nil {
		return err
	}

	if err := m.broker.SendMessages(ctx, ch, tablename); err != nil {
		return err
	}

	return nil
}

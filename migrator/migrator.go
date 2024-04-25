package migrator

import (
	"database/sql"
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
	donor *sql.DB
	recipient *sql.DB
	broker *broker
}

type Migrator interface {
	Migrate()
}

func Migrate(
	donorConfig, recipientConfig DBConfig,
	kafkaConfig KafkaConfig,
) error {
	m, err := newMigrator(&donorConfig, &recipientConfig, &kafkaConfig)
	if err != nil {
		return err
	}
	m.Start()
	defer m.Stop()

	return nil
}

func newMigrator(
	donorConfig, recipientConfig *DBConfig,
	kafkaConfig *KafkaConfig,
) (*migrator, error) {
	donor, err := newDatabase(donorConfig)
	if err != nil {
		return nil, err
	}
	recipient, err := newDatabase(recipientConfig)
	if err != nil {
		return nil, err
	}
	broker, err := newBroker(kafkaConfig)
	if err != nil {
		return nil, err
	}

	return &migrator{
		donor: donor,
		recipient: recipient,
		broker: broker,
	}, nil
}

func (m *migrator) Start() {
	m.broker.Run()
}

func (m *migrator) Stop() {
	m.broker.Stop()
	m.donor.Close()
	m.recipient.Close()
}

func (m *migrator) ReadFromDB() error {
	tx, err := m.donor.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Commit()
	return nil
}

func (m *migrator) InsertToDB() error {
	tx, err := m.recipient.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Commit()
	return nil
}

func (m *migrator) SendToKafka() {
	m.broker.SendMessages()
}

func (m *migrator) GetFromKafka() {
	m.broker.GetMessages()
}

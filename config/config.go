package config

type Driver string

const (
	Postgres Driver = "postgres"
	SQLite   Driver = "sqlite"
	MySQL    Driver = "mysql"
)

type Config struct {
	Database        DatabaseConfig `yaml:"database" env_required:"true"`
	Kafka           KafkaConfig    `yaml:"kafka" env_required:"true"`
	BatchSize       int32          `yaml:"batch_size" env_required:"false"`
	OutputTablename string         `yaml:"output_tablename" env_required:"true"`
	MigrationPath   string         `yaml:"migration_path" env_required:"false"`
	TestCountRows   int            `yaml:"test_count_rows" env_required:"false"`
}

type DatabaseConfig struct {
	Driver Driver `yaml:"driver"`
	DBUri  string `yaml:"db_uri"`
}

type KafkaConfig struct {
	URI  string `yaml:"broker_uri"`
	Topic string `yaml:"topic"`
}

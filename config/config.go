package config

type Driver string

const (
	Postgres Driver = "postgres"
	SQLite   Driver = "sqlite"
	MySQLite Driver = "mysqlite"
)

type Config struct {
	Database      DatabaseConfig `yaml:"database" env_required:"true"`
	Kafka         KafkaConfig    `yaml:"kafka" env_required:"true"`
	BatchSize     int32            `yaml:"batch_size" env_required:"false"`
	Tables        []string       `yaml:"tables" env_required:"false"`
	MigrationPath string         `yaml:"migration_path" env_required:"false"`
	TestCountRows int            `yaml:"test_count_rows" env_required:"false"`
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

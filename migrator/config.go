package migrator

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
### Create config

```
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
```

### Create migrator
```
m, err := migrator.New(dbcfg, kafkacfg)
if err != nil {
	log.Fatalf("Failed to create migrator: %v\n", err)
}
defer m.Close()
```

### Add DBName
```
db := "posrgres"
```

### Run migrate
```
err = m.Migrate(db)
if err != nil {
	log.Printf("Failed to migrate db `%s`: %v\n", db, err)
} else {
	log.Printf("Suscessfully migrated `%s`", db)
}
```

### DO NOT FORGET EDIT `./migrator/kafka/getMsgsFromRows`. Add your structs to switch/case

### MAKEFILE

#### up docker 
```
make docker
```

#### migrate test data
```
make migrate
```

#### run
```
make run
```
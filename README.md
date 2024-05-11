### Add config.yaml to ./config path
 


```
database: 
  driver: "postgres"
  user: "postgres"
  password: "postgres"
  host: "localhost"
  port: 5432
  db_name: "postgres"            //donor db name
kafka:
  host: "localhost"
  port: 29092
batch_size: 10000
tables: ["test", "data", "donor", "users", "drop table users", "users;drop table users; users"]
```

For tests add to config file:
```
migration_path: "./migrations"   //path to the migration sql files
test_count_rows: 25000           //count of test rows

```

### Create migrator
```
cfg := confighelper.MustLoadConfig()  //load config

m, err := migrator.New(
		&cfg.Database,
		&cfg.Kafka, 
		cfg.BatchSize,
	)
if err != nil {
	log.Fatalf("Failed to create migrator: %v\n", err)
}
defer m.Close()
```

### Run migrate loop
```
for _, table := range cfg.Tables {
		err = m.Migrate(table)
		if err != nil {
			log.Printf("Failed to migrate table `%s`: %v\n", table, err)
		} else {
			log.Printf("Successful migrated table `%s`", table)
		}
	}
```

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

#### test
```
make test
```

#### linter
```
make lint
```
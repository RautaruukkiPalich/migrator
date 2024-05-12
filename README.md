### Add config.yaml to ./config path

```
database: 
  driver: "postgres"
  db_uri: "postgres://postgres:postgres@localhost:5441/donor?sslmode=disable"
kafka:
  broker_uri: "localhost:29092"
  topic: "migrator"
batch_size: 10000
```

For tests add to config file:
```
migration_path: "./migrations"   //path to the migration sql files
test_count_rows: 25000           //count of test rows

```

### Create migrator
```
cfg := confighelper.MustLoadConfig()  //load config

m, err := migrator.New(cfg)
if err != nil {
	log.Fatalf("Failed to create migrator: %v\n", err)
}
defer m.Close()
```

### Edit querySelect.sql file in `./resourses` directory
#### for example:
```
SELECT * FROM donor;
```

#### or
```
SELECT Orders.OrderID, Customers.CustomerName, Orders.OrderDate
FROM Orders
INNER JOIN Customers ON Orders.CustomerID=Customers.CustomerID;
```

### Run migrate
```
err = m.Migrate(tablename)
if err != nil {
	log.Printf("Failed to migrate table `%s`: %v\n", table, err)
} else {
	log.Printf("Successful migrated table `%s`", table)
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

#### linter
```
make lint
```
package database

import (
	"context"
	"database/sql"
	"time"

	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"

	"github.com/facily-tech/go-core/env"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	// DBPrefix is the prefix for all environment variables related to the database.
	DBPrefix = "DB_"
	// PostgresDriverName is the name of the postgres driver.
	PostgresDriverName = "pgx"
	// Postgres is the enum for postgres database.
	Postgres database = "postgres"
)

// Drivers is a map with Database enum and its driver name.
var drivers = map[database]string{
	Postgres: PostgresDriverName,
}

type database string

type config struct {
	DSN                  string        `env:"DSN,required"`
	DSNTest              string        `env:"DSN_TEST"`
	MaxOpenConn          int           `env:"MAX_OPEN,default=10"`
	MaxIdleTime          time.Duration `env:"MAX_IDLE_DURATION,default=3m"`
	TracerDatadogEnabled bool          `env:"TRACER_DATADOG_ENABLED,default=false"`
	TracerServiceName    string        `env:"TRACE_SERVICE_NAME,default=database"`
}

// InitDB initializes a new database connection.
func InitDB(database database) (*gorm.DB, *sql.DB, error) {
	var dbConfig config
	if err := env.LoadEnv(context.Background(), &dbConfig, DBPrefix); err != nil {
		return nil, nil, errors.Wrap(err, "cannot load db environment variable")
	}

	return open(database, dbConfig)
}

// open opens a new database connection.
func open(database database, config config) (*gorm.DB, *sql.DB, error) {

	datasource := config.DSN
	if len(config.DSNTest) > 0 {
		datasource = config.DSNTest
	}

	var err error
	var sqlDB *sql.DB

	driverName := drivers[database]

	if config.TracerDatadogEnabled {
		sqltrace.Register(driverName, &stdlib.Driver{}, sqltrace.WithServiceName(config.TracerServiceName))
		sqlDB, err = sqltrace.Open(driverName, datasource)
		if err != nil {
			return nil, nil, err
		}
	} else {
		sqlDB, err = sql.Open(driverName, datasource)
		if err != nil {
			return nil, nil, err
		}
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConn)
	sqlDB.SetConnMaxIdleTime(config.MaxIdleTime)

	db, err := gormtrace.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return db, sqlDB, nil
}

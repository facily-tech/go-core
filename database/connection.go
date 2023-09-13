package database

import (
	"context"
	"database/sql"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	mongotrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go.mongodb.org/mongo-driver/mongo"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"

	"github.com/facily-tech/go-core/env"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	// DBPrefix is the prefix for all environment variables related to the database.
	DBPrefix = "DB_"
	// PostgresDriverName is the name of the postgres driver.
	PostgresDriverName = "pgx"
	// MySQLDriverName is the name of the mysql driver.
	MySQLDriverName = "mysql"
	// Postgres is the enum for postgres database.
	Postgres database = "postgres"
	// MySQL is the enum for mysql database.
	MySQL database = "mysql"
)

// Drivers is a map with Database enum and its driver name.
var drivers = map[database]string{
	Postgres: PostgresDriverName,
	MySQL:    MySQLDriverName,
}

type database string

type config struct {
	DSN                  string        `env:"DSN,required"`
	DSNTest              string        `env:"DSN_TEST"`
	MaxOpenConn          int           `env:"MAX_OPEN,default=10"`
	MaxIdleTime          time.Duration `env:"MAX_IDLE_DURATION,default=1m"`
	MaxLifetime          time.Duration `env:"MAX_LIFETIME_DURATION,default=5m"`
	TracerDatadogEnabled bool          `env:"TRACER_DATADOG_ENABLED,default=false"`
}

// InitDB initializes a new database connection.
func InitDB(database database) (*gorm.DB, *sql.DB, error) {
	return initDB(database, DBPrefix)
}

// InitDBWithPrefix initializes a new database connection with a prefix.
func InitDBWithPrefix(database database, dbPrefix string) (*gorm.DB, *sql.DB, error) {
	return initDB(database, dbPrefix)
}

// InitMongoDB initializes a new mongo database connection.
func InitMongoDB() (*mongo.Client, error) {
	return initMongoDB(DBPrefix)
}

// InitMongoDBWithPrefix initializes a new mongo database connection with a prefix.
func InitMongoDBWithPrefix(dbPrefix string) (*mongo.Client, error) {
	return initMongoDB(dbPrefix)
}

func loadEnv(dbPrefix string) (*config, error) {
	var dbConfig config
	if err := env.LoadEnv(context.Background(), &dbConfig, dbPrefix); err != nil {
		return nil, errors.Wrap(err, "cannot load db environment variable")
	}

	return &dbConfig, nil
}

// InitMongoDB initializes a new mongo database connection.
func initMongoDB(dbPrefix string) (*mongo.Client, error) {
	dbConfig, err := loadEnv(dbPrefix)
	if err != nil {
		return nil, err
	}

	return openMongoConn(dbConfig)
}

// openMongoConn opens a new mongo database connection.
func openMongoConn(dbConfig *config) (*mongo.Client, error) {
	opts := options.Client()
	opts.Monitor = mongotrace.NewMonitor()
	opts.ApplyURI(dbConfig.DSN)
	opts.MaxConnIdleTime = &dbConfig.MaxIdleTime

	maxPoolSize := uint64(dbConfig.MaxOpenConn)
	opts.MaxPoolSize = &maxPoolSize

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open mongo connection")
	}

	return client, nil
}

// InitDB initializes a new database connection.
func initDB(database database, dbPrefix string) (*gorm.DB, *sql.DB, error) {
	dbConfig, err := loadEnv(dbPrefix)
	if err != nil {
		return nil, nil, err
	}

	return open(database, dbConfig)
}

// open opens a new database connection.
func open(database database, config *config) (*gorm.DB, *sql.DB, error) {
	datasource := config.DSN
	if len(config.DSNTest) > 0 {
		datasource = config.DSNTest
	}

	var err error
	var sqlDB *sql.DB

	driverName := drivers[database]

	if config.TracerDatadogEnabled {
		if database == Postgres {
			sqltrace.Register(driverName, &stdlib.Driver{})
		} else if database == MySQL {
			sqltrace.Register(driverName, &mysqldriver.MySQLDriver{})
		}
	}

	sqlDB, err = sql.Open(driverName, datasource)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot open database connection")
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConn)
	sqlDB.SetConnMaxIdleTime(config.MaxIdleTime)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	var dialector gorm.Dialector
	if database == Postgres {
		dialector = postgres.New(postgres.Config{Conn: sqlDB})
	} else if database == MySQL {
		dialector = mysql.New(mysql.Config{Conn: sqlDB})
	}

	db, err := gormtrace.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot open a gorm connection")
	}

	return db, sqlDB, nil
}

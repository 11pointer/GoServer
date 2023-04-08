package database

import (
	"11pointer/logger"
	"context"
	"contrib.go.opencensus.io/integrations/ocsql"
	"database/sql/driver"
	"fmt"
	entsql "github.com/facebook/ent/dialect/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/mysql"
	"go.uber.org/zap"
	"os"
	"time"
)

var (
	TemporalHostPort  = "localhost:7233"
	TemporalNamespace = "default"
)

type connector struct {
	dsn string
}

func (c connector) Connect(context.Context) (driver.Conn, error) {
	return c.Driver().Open(c.dsn)
}

func (connector) Driver() driver.Driver {
	return ocsql.Wrap(
		mysql.MySQLDriver{},
		ocsql.WithAllTraceOptions(),
		ocsql.WithRowsClose(false),
		ocsql.WithRowsNext(false),
		ocsql.WithDisableErrSkip(true),
	)
}

func init() {
	configDir := os.Getenv("CONFIG_DIR")
	//if configDir == "" {
	//	configDir = "./config"
	//}
	viper.AddConfigPath(configDir)
	viper.SetConfigType("json")
	viper.SetConfigFile(configDir + "/config.json")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		logger.LOG.Panic("VIPER config read error", zap.Error(err))
	}
}

func GetDBDriver() (*entsql.Driver, error) {
	host := viper.GetString("host")
	port := viper.GetInt("port")
	username := viper.GetString("db_username")
	password := viper.GetString("db_password")
	databaseName := viper.GetString("database_name")
	dataSource := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true", username, password, host, port, databaseName)
	logger.LOG.Info("DB Connection", zap.String("host", host), zap.Int("port", port), zap.String("databaseName", databaseName))
	driver, _ := Open(dataSource)
	return driver, nil
}

func Open(dsn string) (*entsql.Driver, error) {

	var err error
	db, err := apmsql.Open("mysql", dsn)
	if err != nil {
		logger.LOG.Panic("Could not open connection to Database", zap.Error(err))
	}
	maxOpenConnections := 25
	maxIdleConnections := 3
	maxConnectionLifeTime := 300 * time.Second
	maxOpenConnections = viper.GetInt("max_open_connections")
	maxIdleConnections = viper.GetInt("max_idle_connections")
	maxConnectionLifeTime = time.Duration(viper.GetInt("max_connection_lifetime_seconds")) * time.Second
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxLifetime(maxConnectionLifeTime)
	err = db.Ping()
	if err != nil {
		logger.LOG.Panic("Could not connect to Database", zap.Error(err))
	} else {
		logger.LOG.Info("Successfully connected to database!")
	}
	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB("mysql", db)
	return drv, nil
}

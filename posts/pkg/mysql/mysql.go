package mysql

import (
	"database/sql"
	"fmt"
	"post_service/config"
	"post_service/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
)

type MySqlDriver struct {
	db     *sql.DB
	config *config.Config
}

var (
	l     logger.Logger
	cfg   *config.Config
	MySql MySqlDriver
)

func init() {
	cfg = config.NewConfig()
	l = *logger.New(cfg.App.LogLevel)
	l.Info("MySQL driver initialized")

	MySql = *NewMySqlDriver(cfg)
	if err := MySql.Create(); err != nil {
		l.Error("Failed to initialize MySQL connection: %v", err)
		return
	}
	l.Info("MySQL connected")
}

func NewMySqlDriver(cfg *config.Config) *MySqlDriver {
	if cfg == nil {
		return nil
	}

	return &MySqlDriver{
		config: cfg,
	}
}

func (d *MySqlDriver) Create() error {
	if d.config == nil {
		return fmt.Errorf("configuration is required")
	}

	db, err := sql.Open("mysql", d.config.MySql.User)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	d.db = db
	return nil
}

func (d *MySqlDriver) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if d.db == nil {
		return nil, sql.ErrConnDone
	}
	return d.db.Query(query, args...)
}

func (d *MySqlDriver) QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	if d.db == nil {
		return nil, sql.ErrConnDone
	}
	return d.db.QueryRow(query, args...), nil
}

func (d *MySqlDriver) Exec(query string, args ...interface{}) (sql.Result, error) {
	if d.db == nil {
		return nil, sql.ErrConnDone
	}
	return d.db.Exec(query, args...)
}

func (d *MySqlDriver) Prepare(query string) (*sql.Stmt, error) {
	if d.db == nil {
		return nil, sql.ErrConnDone
	}
	return d.db.Prepare(query)
}

func (d *MySqlDriver) Begin() (*sql.Tx, error) {
	if d.db == nil {
		return nil, sql.ErrConnDone
	}
	return d.db.Begin()
}

func (d *MySqlDriver) Close() error {
	if d.db == nil {
		return nil
	}
	return d.db.Close()
}

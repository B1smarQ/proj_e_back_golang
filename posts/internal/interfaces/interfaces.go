package interfaces

import (
	"database/sql"
)

type DatabaseDriver interface {
	Create() error
	Close() error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) (*sql.Row, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Begin() (*sql.Tx, error)
}

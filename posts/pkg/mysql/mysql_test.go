package mysql

import (
	"database/sql"
	"testing"

	"post_service/config"
)

// mockDB is a mock implementation of sql.DB
type mockDB struct {
	*sql.DB
}

func TestMySqlDriver_Create(t *testing.T) {
	// Test creating driver without configuration
	driver := NewMySqlDriver(nil)
	if driver != nil {
		t.Error("expected nil driver when creating without configuration")
	}

	// Test creating driver with configuration
	cfg := &config.Config{
		MySQL: config.MySQLConfig{
			Host:     "localhost",
			Port:     "3306",
			User:     "test",
			Password: "test",
			Database: "test",
		},
		App: config.AppConfig{
			LogLevel: "info",
		},
	}

	driver = NewMySqlDriver(cfg)
	if driver == nil {
		t.Error("expected non-nil driver when creating with configuration")
	}
}

func TestMySqlDriver_Methods(t *testing.T) {
	cfg := &config.Config{
		MySQL: config.MySQLConfig{
			Host:     "localhost",
			Port:     "3306",
			User:     "test",
			Password: "test",
			Database: "test",
		},
		App: config.AppConfig{
			LogLevel: "info",
		},
	}

	driver := NewMySqlDriver(cfg)
	if driver == nil {
		t.Fatal("failed to create driver")
	}

	// Test Query
	_, err := driver.Query("SELECT 1")
	if err == nil {
		t.Error("expected error when querying without connection")
	}

	// Test QueryRow
	_, err = driver.QueryRow("SELECT 1")
	if err == nil {
		t.Error("expected error when querying row without connection")
	}

	// Test Exec
	_, err = driver.Exec("INSERT INTO test VALUES (1)")
	if err == nil {
		t.Error("expected error when executing without connection")
	}

	// Test Prepare
	_, err = driver.Prepare("SELECT 1")
	if err == nil {
		t.Error("expected error when preparing statement without connection")
	}

	// Test Begin
	_, err = driver.Begin()
	if err == nil {
		t.Error("expected error when beginning transaction without connection")
	}

	// Test Close
	err = driver.Close()
	if err != nil {
		t.Errorf("unexpected error when closing driver: %v", err)
	}
}

func TestMySqlDriver_Query(t *testing.T) {
	driver := &MySqlDriver{}
	_, err := driver.Query("SELECT * FROM test")
	if err != sql.ErrConnDone {
		t.Errorf("expected error %v, got %v", sql.ErrConnDone, err)
	}
}

func TestMySqlDriver_QueryRow(t *testing.T) {
	driver := &MySqlDriver{}
	_, err := driver.QueryRow("SELECT * FROM test")
	if err != sql.ErrConnDone {
		t.Errorf("expected error %v, got %v", sql.ErrConnDone, err)
	}
}

func TestMySqlDriver_Exec(t *testing.T) {
	driver := &MySqlDriver{}
	_, err := driver.Exec("INSERT INTO test VALUES (?)", "test")
	if err != sql.ErrConnDone {
		t.Errorf("expected error %v, got %v", sql.ErrConnDone, err)
	}
}

func TestMySqlDriver_Prepare(t *testing.T) {
	driver := &MySqlDriver{}
	_, err := driver.Prepare("SELECT * FROM test")
	if err != sql.ErrConnDone {
		t.Errorf("expected error %v, got %v", sql.ErrConnDone, err)
	}
}

func TestMySqlDriver_Begin(t *testing.T) {
	driver := &MySqlDriver{}
	_, err := driver.Begin()
	if err != sql.ErrConnDone {
		t.Errorf("expected error %v, got %v", sql.ErrConnDone, err)
	}
}

func TestMySqlDriver_Close(t *testing.T) {
	driver := &MySqlDriver{}
	err := driver.Close()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

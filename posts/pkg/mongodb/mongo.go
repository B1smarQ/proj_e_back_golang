package mongodb

import (
	"context"
	"database/sql"
	"errors"
	"post_service/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type MongoDriver struct {
	client *mongo.Client
}

func (m *MongoDriver) Create(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		return err
	}
	m.client = client
	return nil
}

func (Mongo MongoDriver) Close() error {
	return Mongo.client.Disconnect(context.TODO())
}

func (Mongo MongoDriver) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, ErrNotImplemented
}

func (Mongo MongoDriver) QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	return nil, ErrNotImplemented
}

func (Mongo MongoDriver) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, ErrNotImplemented
}

func (Mongo MongoDriver) Prepare(query string) (*sql.Stmt, error) {
	return nil, ErrNotImplemented
}

func (Mongo MongoDriver) Begin() (*sql.Tx, error) {
	return nil, ErrNotImplemented
}

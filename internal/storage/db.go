package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	conn *mongo.Client
	db   *mongo.Database
}

func NewMongo(config *Config) (*Database, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.MongoURI).SetServerAPIOptions(serverAPI)
	conn, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	db := conn.Database(config.MongoDatabase)
	return &Database{
		conn: conn,
		db:   db,
	}, nil
}

func (db *Database) Close() error {
	err := db.conn.Disconnect(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetCollection(name string) *mongo.Collection {
	return db.db.Collection(name)
}

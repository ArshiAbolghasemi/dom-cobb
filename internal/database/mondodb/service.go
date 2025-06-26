package mondodb

import (
	"context"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	once   sync.Once
)

func GetClient() *mongo.Client {
	once.Do(func() {
		uri, err := GetURI()
		if err != nil {
			panic("Failed to get MongoDB URI: " + err.Error())
		}
		connTimeout, err := GetConnectionTimeout()
		if err != nil {
			panic("Failed to get MongoDB connection timeout: " + err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), connTimeout*time.Second)
		defer cancel()

		conn, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			panic("Failed to connect to MongoDB: " + err.Error())
		}

		if err := conn.Ping(ctx, nil); err != nil {
			panic("Failed to ping MongoDB: " + err.Error())
		}

		client = conn
	})
	return client
}

func GetDB() *mongo.Database {
	client := GetClient()

	dbname, exists := os.LookupEnv("MONGO_DBNAME")
	if !exists {
		panic("Mongo dbname is undefined")
	}

	return client.Database(dbname)
}

func GetCollection(collection string) *mongo.Collection {
	return GetDB().Collection(collection)
}

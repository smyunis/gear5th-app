package mongodbpersistence

import (
	"context"
	"fmt"

	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBStore interface {
	Database() *mongo.Database
}

var client *mongo.Client

type MongoDBStoreBootstrap struct {
	config infrastructure.ConfigurationProvider
}

func NewMongoDBStoreBootstrap(config infrastructure.ConfigurationProvider) MongoDBStoreBootstrap {
	bootstrap := MongoDBStoreBootstrap{
		config,
	}
	bootstrap.initDB()
	return bootstrap
}

func (m MongoDBStoreBootstrap) initDB() error {

	if client == nil {
		connString := m.config.Get("MONGODB_URL", "")
		clientOptions := options.Client().ApplyURI(connString)
		var err error
		client, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			panic(fmt.Errorf("could not connect to mongo db database: %w", err))
		}
	}
	return nil
}

func (MongoDBStoreBootstrap) Database() *mongo.Database {
	databaseName := "gear5th"
	return client.Database(databaseName)
}

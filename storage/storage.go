package storage

import (
	"context"
	"os"
	"strings"
	"time"

	"slackbot-test/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName = "nerdcoin"
	collectionName = "slack-messages"
)

var Db *mongo.Collection
var Ctx context.Context
func Setup() {
	Ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbClient, err := mongo.Connect(Ctx, options.Client().ApplyURI(os.Getenv("MONGO")))
	if err != nil {
		logger.Error(err.Error())
	}
	defer func() {
		if err = dbClient.Disconnect(Ctx); err != nil {
			panic(err)
		}
	}()

	dbNames, err := dbClient.ListDatabaseNames(Ctx, bson.D{})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	var dbExists = false;
	for _, db := range dbNames {
		if strings.Compare(strings.ToLower(db), dbName) == 0 {
			dbExists = true
			logger.Info(dbName + " database exists", )
			break
		}
	}

	if !dbExists {
		logger.Info(dbName + " database does not exists. Creating database...", )
		err := dbClient.Database(dbName).CreateCollection(Ctx, collectionName)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		logger.Info(dbName + " database with " + collectionName + " collection created successfully.")
	} else {
		var collectionExists = false;
		collectionNames, err := dbClient.Database(dbName).ListCollectionNames(Ctx, bson.D{})
		if err != nil {
			logger.Error(err.Error())
			return
		}
		for _, collection := range collectionNames {
			if strings.Compare(strings.ToLower(collection), collectionName) == 0 {
				collectionExists = true
				logger.Info(collectionName + " collection exists")
				break
			}
		}
		if !collectionExists {
			logger.Info(collectionName + " collection does not exists. Creating collection...")
			err := dbClient.Database(dbName).CreateCollection(Ctx, collectionName)
			if err != nil {
				logger.Error(err.Error())
				return
			}
			logger.Info(collectionName + " collection created successfully")
		}

		Db = dbClient.Database(dbName).Collection(collectionName)
		// count, err := collection.CountDocuments(Ctx, bson.D{})
		// logger.Info("This many documents: " + string(count))
	}
}


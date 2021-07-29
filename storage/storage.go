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
	DbName         = "nerdcoin"
	CollectionName = "slack-messages"
)

func Setup() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO")))
	if err != nil {
		logger.Error(err.Error())
	}
	defer func() {
		if err = dbClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	dbNames, err := dbClient.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	var dbExists = false
	for _, db := range dbNames {
		if strings.Compare(strings.ToLower(db), DbName) == 0 {
			dbExists = true
			logger.Info(DbName + " database exists", )
			break
		}
	}

	if !dbExists {
		logger.Info(DbName + " database does not exists. Creating database...", )
		err := dbClient.Database(DbName).CreateCollection(ctx, CollectionName)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		logger.Info(DbName + " database with " + CollectionName + " collection created successfully.")
	} else {
		var collectionExists = false
		collectionNames, err := dbClient.Database(DbName).ListCollectionNames(ctx, bson.D{})
		if err != nil {
			logger.Error(err.Error())
			return
		}
		for _, collection := range collectionNames {
			if strings.Compare(strings.ToLower(collection), CollectionName) == 0 {
				collectionExists = true
				logger.Info(CollectionName + " collection exists")
				break
			}
		}
		if !collectionExists {
			logger.Info(CollectionName + " collection does not exists. Creating collection...")
			err := dbClient.Database(DbName).CreateCollection(ctx, CollectionName)
			if err != nil {
				logger.Error(err.Error())
				return
			}
			logger.Info(CollectionName + " collection created successfully")
		}

		collection := dbClient.Database(DbName).Collection(CollectionName)
		count, err := collection.CountDocuments(ctx, bson.D{})
		logger.Info("This many documents: " + string(count))
	}
}


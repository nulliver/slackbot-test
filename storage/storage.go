package storage

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"slackbot-test/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName         = "nerdcoin"
	collectionName = "slack-messages"
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
		if strings.Compare(strings.ToLower(db), dbName) == 0 {
			dbExists = true
			logger.Info(dbName + " database exists", )
			break
		}
	}

	if !dbExists {
		logger.Info(dbName + " database does not exists. Creating database...", )
		err := dbClient.Database(dbName).CreateCollection(ctx, collectionName)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		logger.Info(dbName + " database with " + collectionName + " collection created successfully.")
	} else {
		var collectionExists = false
		collectionNames, err := dbClient.Database(dbName).ListCollectionNames(ctx, bson.D{})
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
			err := dbClient.Database(dbName).CreateCollection(ctx, collectionName)
			if err != nil {
				logger.Error(err.Error())
				return
			}
			logger.Info(collectionName + " collection created successfully")
		}

		collection := dbClient.Database(dbName).Collection(collectionName)
		count, err := collection.CountDocuments(ctx, bson.D{})
		logger.Info("This many documents: " + string(count))
	}
}

func SaveTransaction(username, message string, usersWithPlusPlus []string) {
	var docs []interface{}

	for _, u := range usersWithPlusPlus {
		doc := bson.D{
			{Key: "fromUser", Value: username},
			{Key: "toUser", Value: u},
			{Key: "message", Value: message},
			{Key: "type", Value: "karma"},
			{Key: "timestamp", Value: time.Now()},
		}
		docs = append(docs, doc)
	}

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

	res, err := dbClient.Database(dbName).Collection(collectionName).InsertMany(ctx, docs)
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Println(res)
}
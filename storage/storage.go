package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
		log.Print(err)
	}
	defer func() {
		if err = dbClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	dbNames, err := dbClient.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		log.Print(err)
		return
	}

	var dbExists = false
	for _, db := range dbNames {
		if strings.Compare(strings.ToLower(db), dbName) == 0 {
			dbExists = true
			log.Printf("'%s' database exists", dbName)
			break
		}
	}

	if !dbExists {
		log.Printf( "'%s' database does not exists. Creating database...", dbName)
		err := dbClient.Database(dbName).CreateCollection(ctx, collectionName)
		if err != nil {
			log.Print(err)
			return
		}
		log.Printf( "'%s' database with '%s' collection created successfully", dbName, collectionName)
	} else {
		var collectionExists = false
		collectionNames, err := dbClient.Database(dbName).ListCollectionNames(ctx, bson.D{})
		if err != nil {
			log.Print(err)
			return
		}
		for _, collection := range collectionNames {
			if strings.Compare(strings.ToLower(collection), collectionName) == 0 {
				collectionExists = true
				log.Printf( "'%s' collection exists", collectionName)
				break
			}
		}
		if !collectionExists {
			log.Printf(  "'%s' collection does not exists. Creating collection...", collectionName)
			err := dbClient.Database(dbName).CreateCollection(ctx, collectionName)
			if err != nil {
				log.Print(err)
				return
			}
			log.Printf("'%s' collection created successfully", collectionName)
		}
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
		log.Printf(err.Error())
	}
	defer func() {
		if err = dbClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	res, err := dbClient.Database(dbName).Collection(collectionName).InsertMany(ctx, docs)
	if err != nil {
		log.Printf(err.Error())
	}
	fmt.Println(res)
}
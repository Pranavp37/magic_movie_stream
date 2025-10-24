package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func MongoDBConnection(uri string) {

	clientoption := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientoption)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to MongoDB")

	mongoClient = client
}

func GetMongoCollection(databasemname, collectionname string) *mongo.Collection {

	if mongoClient == nil {
		log.Fatalf("Failed to get mongpo client try to cnnect database")
	}
	return mongoClient.Database(databasemname).Collection(collectionname)
}

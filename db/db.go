package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	MongoClient   *mongo.Client
	MongoDatabase *mongo.Database
)

func ConnectMongo(mongoUri, mongoDatabaseName string) {
	mongoClient, err := mongo.Connect(nil, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Panicf("Error connecting to MongoDB: %v", err)
	}

	log.Print("Successfully connected to MongoDB")

	MongoClient = mongoClient
	MongoDatabase = mongoClient.Database(mongoDatabaseName)
}

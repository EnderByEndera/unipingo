package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection struct {
	Client  *mongo.Client
	Context context.Context
}

var mongoConn *MongoConnection

func InitMongo() *MongoConnection {
	cli, err := mongo.NewClient(options.Client().
		ApplyURI("mongodb://localhost:10200").
		SetTimeout(10 * time.Second))
	if err != nil {
		log.Fatal(err)
	}
	err = cli.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	mongoConnection := MongoConnection{Client: cli}

	return &mongoConnection
}

func GetMongoConn() *MongoConnection {
	if mongoConn == nil {
		mongoConn = InitMongo()
	}
	return mongoConn
}

func GetCollection(collectionName string) (collection *mongo.Collection) {
	conn := GetMongoConn()
	collection = conn.Client.Database("blog").Collection(collectionName)
	return
}

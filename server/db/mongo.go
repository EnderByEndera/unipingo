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
	cli, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:10200"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = cli.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	mongoConnection := MongoConnection{Client: cli, Context: ctx}

	return &mongoConnection
}

func GetMongoConn() *MongoConnection {
	if mongoConn == nil {
		mongoConn = InitMongo()
	}
	return mongoConn
}

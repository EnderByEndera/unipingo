package main

import (
	"context"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"testing"

	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestBatchSize(t *testing.T) {
	opts := options.Find().SetLimit(20)
	cur, err := db.GetCollection("majors").Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		panic(err)
	}
	var results []*models.Major
	for cur.Next(context.TODO()) {
		_, err := cur.Current.Elements()
		if err != nil {
			panic(err)
		}
	}
	assert.Equal(t, len(results), 20)
}

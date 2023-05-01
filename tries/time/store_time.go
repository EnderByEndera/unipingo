package main

import (
	"context"
	"fmt"
	"melodie-site/server/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Time struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Time time.Time          `bson:"time" json:"time"`
}

func main() {
	t := time.Now()
	fmt.Println(t.Add(-10 * time.Minute).Unix())
	cur, err := db.GetCollection("time").Find(context.TODO(), bson.M{
		"time": bson.M{
			"$gt": t.Add(-11 * time.Minute),
		},
	})
	if err != nil {
		panic(err)
	}
	var result []*Time
	err = cur.All(context.TODO(), &result)
	if err != nil {
		panic(err)
	}
	for _, res := range result {
		fmt.Println(res.ID, res.Time.Unix())
	}
}

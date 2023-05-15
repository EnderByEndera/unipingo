package main

import (
	"context"
	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"testing"
)

func TestRollBack(t *testing.T) {
	question := new(models.QuestionBoxQuestion)
	err := db.GetCollection("questions").FindOne(context.TODO(), bson.M{"title": "My Question"}).Decode(question)
	assert.Equal(t, err, nil)

	_, err = db.GetMongoConn().UseSession(nil, func(sessionContext mongo.SessionContext) (result interface{}, err error) {
		err = db.GetCollection("questions").FindOneAndUpdate(sessionContext,
			bson.M{"_id": question.ID},
			bson.M{"$set": bson.M{"title": "My New Question"}}).Err()
		if err != nil {
			return
		}

		err = db.GetCollection("questions").FindOneAndUpdate(sessionContext,
			bson.M{"title": "My Question"},
			bson.M{"$set": bson.M{"name": "Ender"}}).Err()
		result = db.GetCollection("questions").FindOne(context.TODO(), bson.M{"title": "My New Question"})
		return
	})

	assert.NotEqual(t, err, nil)
	err = db.GetCollection("questions").FindOne(context.TODO(), bson.M{"title": "My Question"}).Decode(question)
	assert.Equal(t, err, nil)
	assert.Equal(t, question.Title, "My Question")
}

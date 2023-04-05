package tests

import (
	"melodie-site/server/utils"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDInterface struct {
	UserID   primitive.ObjectID
	AnswerID primitive.ObjectID
}

func TestMutex(t *testing.T) {
	mtx := utils.KeyedMutex{}
	key1 := IDInterface{UserID: primitive.NewObjectID(), AnswerID: primitive.NewObjectID()}
	// key2 := IDInterface{UserID: primitive.NewObjectID(), AnswerID: primitive.NewObjectID()}
	mtx.Lock(key1)
	// mtx.Lock(key1)
	mtx.Unlock(key1)
}

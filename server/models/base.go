package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EntityWithName struct {
	ID   primitive.ObjectID `json:"id" bson:"id"`
	Name string             `json:"name" bson:"name"`
}

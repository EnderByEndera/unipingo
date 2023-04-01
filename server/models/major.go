package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Major struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Code        string             `bson:"code" json:"code"`
	Name        string             `bson:"name" json:"name"`
	Category    string             `bson:"category" json:"category"`
	Description string             `bson:"description" json:"description"`
	Topics      []string           `bson:"topics" json:"topics"`
}

package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Major struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Code        string             `bson:"code" json:"code"`
	Name        string             `bson:"name" json:"name"`
	Category    string             `bson:"category" json:"category"` // 大类
	FirstLevel  string             `bson:"firstLevel" json:"firstLevel"` // 所属一级学科
	Description string             `bson:"description" json:"description"`
	Topics      []string           `bson:"topics" json:"topics"`
}

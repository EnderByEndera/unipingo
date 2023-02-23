package models

import (
	"encoding/json"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	LIKE            uint8 = 0
	DISLIKE         uint8 = 1
	ADD_TO_FAVORITE uint8 = 2

	TYPE_ARTICLE uint8 = 0
	TYPE_COMMENT uint8 = 1
	TYPE_REPLY   uint8 = 2
)

type UserAction struct {
	Type   uint8 `bson:"type" json:"type"`
	UserID int   `bson:"userID" json:"userID"`
	Time   int64 `bson:"time" json:"time"`
}

type Post struct {
	UUID        uuid.UUID          `bson:"uuid" json:"uuid"`
	UserID      int                `bson:"userID" json:"userID"`
	Content     string             `bson:"content" json:"content"`
	UserActions []UserAction       `bson:"userActions" json:"userActions"`
	DocumentID  primitive.ObjectID `bson:"_id,omitempty" json:"_oid"`
	Title       string             `bson:"title" json:"title"`
	Comments    []Comment          `bson:"comments" json:"comments"`
}

func (content *Post) ToIndentedJSON() string {
	obj, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(obj)
}

// 评论
type Comment struct {
	UUID        uuid.UUID    `bson:"uuid" json:"uuid"`
	UserID      int          `bson:"userID" json:"userID"`
	Content     string       `bson:"content" json:"content"`
	UserActions []UserAction `bson:"userActions" json:"userActions"`
	Replies     []Reply      `bson:"replies" json:"replies"`
}

// 回复
type Reply struct {
	UUID        uuid.UUID    `bson:"uuid" json:"uuid"`
	UserID      int          `bson:"userID" json:"userID"`
	Content     string       `bson:"content" json:"content"`
	UserActions []UserAction `bson:"userActions" json:"userActions"`
	ToUUID      uuid.UUID    `bson:"toUUID" json:"toUUID"`
	// Replies []Reply `bson:"replies" json:"replies"`
}

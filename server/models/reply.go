package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReplyStats struct {
	Likes    int `bson:"likes" json:"likes"`
	Dislikes int `bson:"dislikes" json:"dislikes"`
}
type Reply struct {
	CreateTime uint64               `bson:"createTime" json:"createTime"`
	UpdateTime uint64               `bson:"updateTime" json:"updateTime"`
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"_oid"`
	UserID     primitive.ObjectID   `bson:"userID" json:"userID"`
	CommentID  primitive.ObjectID   `bson:"commentID" json:"commentID"`
	Statistics ReplyStats           `bson:"statistics" json:"statistics"`
	Liked      []primitive.ObjectID `bson:"liked" json:"liked"` // 点赞的用户
	Content    string               `bson:"content" json:"content"`
}

package models

import (
	"melodie-site/server/utils"

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

type Like struct {
	UserID    primitive.ObjectID `bson:"userID" json:"userID"`
	Position  bool               `bson:"position" json:"position"` //
	TimeStamp int64              `bson:"time" json:"time"`
}

type Favorite struct {
	UserID                 primitive.ObjectID `bson:"userID" json:"userID"`
	TimeStamp              int64              `bson:"time" json:"time"`
	FavoriteCollectionUUID string             `bson:"favoriteCollectionUUID" json:"favoriteCollectionUUID"`
}

type SimpleStats struct {
	Likes    int `bson:"likes" json:"likes"`
	Dislikes int `bson:"dislikes" json:"dislikes"`
}

type AnswerStats struct {
	Approves         int `bson:"approves" json:"approves"`
	Disapproves      int `bson:"disapproves" json:"disapproves"`
	AlumnApproves    int `bson:"alumnApproves" json:"alumnApproves"`
	AlumnDisapproves int `bson:"alumnDisapproves" json:"alumnDisapproves"`
	Favorites        int `bson:"favorites" json:"favorites"`
}

type AnswerCategory string

const (
	AnswerAboutHEI   = "HEI"
	AnswerAboutMajor = "Major"
)

type Answer struct {
	CreateTime       uint64               `bson:"createTime" json:"createTime"`
	UpdateTime       uint64               `bson:"updateTime" json:"updateTime"`
	ID               primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID           primitive.ObjectID   `bson:"userID" json:"userID"`
	Category         AnswerCategory       `bson:"category" json:"category"`   // 属于学校或者专业
	BelongsTo        EntityWithName       `bson:"belongsTo" json:"belongsTo"` // 属于什么学校或者专业
	Question         string               `bson:"question" json:"question"`
	Statistics       AnswerStats          `bson:"statistics" json:"statistics"`
	ApprovedUsers    []primitive.ObjectID `bson:"approvedUsers" json:"approvedUsers"`       // 点赞的用户
	DisapprovedUsers []primitive.ObjectID `bson:"disapprovedUsers" json:"disapprovedUsers"` // 点赞的用户
	FavoritedUsers   []primitive.ObjectID `bson:"favoritedUsers" json:"favoritedUsers"`     // 添加到收藏的用户
	Title            string               `bson:"title" json:"title"`
	Content          string               `bson:"content" json:"content"`
}

type NewAnswerRequest struct {
	Category AnswerCategory     `json:"category"`
	Question string             `json:"question"`
	Content  string             `json:"content"`
	EntityID primitive.ObjectID `json:"entityID"` // 学校或者专业的ID
}

type ApproveOrDisapproveAnswerRequest struct {
	AnsID   primitive.ObjectID `json:"ansID"`
	Approve bool               `json:"approve"`
}

func (ans *Answer) Init() {
	ans.ApprovedUsers = make([]primitive.ObjectID, 0)
	ans.DisapprovedUsers = make([]primitive.ObjectID, 0)
	ans.FavoritedUsers = make([]primitive.ObjectID, 0)
}

func (content *Answer) ToIndentedJSON() string {
	return utils.ToIndentedJSON(content)
}

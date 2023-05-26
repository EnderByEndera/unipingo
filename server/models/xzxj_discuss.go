package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type XZXJUser struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`                    // 学长学姐说集合中的_id
	UserID            primitive.ObjectID `bson:"userID" json:"userID"`                       // 参与学长学姐说的用户的真实ID，实际使用应该使用此ID
	Sections          []string           `bson:"sections" json:"sections"`                   // 学长学姐说用户参与的分享计划
	Picture           string             `bson:"picture" json:"picture"`                     // 学长学姐说用户照片
	Motto             string             `bson:"motto" json:"motto"`                         // 格言
	ManagedActivities string             `bson:"managedActivities" json:"managedActivities"` // 主理活动
	Experience        string             `bson:"experience" json:"experience"`               // 个人经历
}

type XZXJUserFormMeta struct {
	XZXJUser
	RealName string   `bson:"realName" json:"realName"`
	UserTags []string `bson:"userTags" json:"userTags"`
}

package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"melodie-site/server/models"
)

type AddOrUpdateXZXJUserReq struct {
	Sections          []string `json:"sections"`          // 学长学姐说用户参与的分享计划
	Picture           string   `json:"picture"`           // 学长学姐说用户的个人照片
	Motto             string   `json:"motto"`             // 学长学姐说用户的座右铭
	ManagedActivities string   `json:"managedActivities"` // 学长学姐说用户的主理活动
	Experience        string   `json:"experience"`        // 学长学姐说用户的个人经历
	RealName          string   `json:"realName"`          // 学长学姐说用户的真实姓名
	UserTags          []string `json:"userTags"`          // 学长学姐说用户的标签
}

type AddOrUpdateXZXJUserRes struct {
	UserID primitive.ObjectID `json:"userID"` // 学长学姐说用户的用户ID
}

type QueryXZXJUserByUserIDReq struct {
}

type QueryXZXJUserByUserIDRes struct {
	XZXJUser *models.XZXJUser `json:"xzxjUser"` // 学长学姐说用户额外信息
}

type DeleteXZXJUserReq struct {
}

type DeleteXZXJUserRes struct {
	Deleted bool `json:"deleted"`
}

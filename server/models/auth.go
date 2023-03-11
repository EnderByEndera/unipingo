package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type roles struct {
	Admin      string // admin可以看到学生资质审核的流程。
	UnpaidUser string
	PaidUser   string
}

var Roles roles = roles{"admin", "unpaid", "paid"}

type WechatInfo struct {
	OpenID  string `json:"openID" bson:"openID"`
	UnionID string `json:"unionID" bson:"unionID"`
}

type User struct {
	OID          primitive.ObjectID `json:"oid" bson:"_id,omitempty"`
	Role         string             `json:"role" bson:"role"`
	Name         string             `json:"name" bson:"name"`
	EMail        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"passwordHash"`
	Avatar       string             `json:"avatar" bson:"avatar"`
	WechatInfo   WechatInfo         `json:"wechatInfo" bson:"wechatInfo"`
}

type UserResponse struct {
	OID        string     `json:"id"`
	Name       string     `json:"name"`
	EMail      string     `json:"email"`
	WechatInfo WechatInfo `json:"wechatInfo"`
}

type UserPublicInfo struct {
	OID    string `json:"id"`
	Avatar string `json:"avatar" bson:"avatar"`
	Name   string `json:"name" bson:"name"`
}

func (userResponse *UserResponse) LoadFromStructUser(user *User) {
	userResponse.OID = user.OID.Hex()
	userResponse.Name = user.Name
	userResponse.EMail = user.EMail
	userResponse.WechatInfo = user.WechatInfo
}

func (user *User) ToPublicInfo() UserPublicInfo {
	return UserPublicInfo{
		OID:    user.OID.Hex(),
		Avatar: user.Avatar,
		Name:   user.Name,
	}
}

type LoginResponse struct {
	UserInfo UserResponse `json:"user"`
	JWTToken string       `json:"jwtToken"`
}

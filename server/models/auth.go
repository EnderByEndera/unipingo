package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WechatInfo struct {
	OpenID  string `json:"openID" bson:"openID"`
	UnionID string `json:"unionID" bson:"unionID"`
}

type User struct {
	OID          primitive.ObjectID `json:"oid" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	EMail        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"passwordHash"`
	WechatInfo   WechatInfo         `json:"wechatInfo" bson:"wechatInfo"`
}

type UserResponse struct {
	OID        string     `json:"id"`
	Name       string     `json:"name"`
	EMail      string     `json:"email"`
	WechatInfo WechatInfo `json:"wechatInfo"`
}

func (userResponse *UserResponse) LoadFromStructUser(user *User) {
	userResponse.OID = user.OID.Hex()
	userResponse.Name = user.Name
	userResponse.EMail = user.EMail
	userResponse.WechatInfo = user.WechatInfo
}

type LoginResponse struct {
	UserInfo UserResponse `json:"user"`
	JWTToken string       `json:"jwtToken"`
}

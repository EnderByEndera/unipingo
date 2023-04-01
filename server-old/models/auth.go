package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RoleAdmin      string = "ROLE_ADMIN"
	RoleUnpaidUser string = "ROLE_UNPAID_USER"
	RolePaidUser   string = "ROLE_PAID_USER"
)

type WechatInfo struct {
	OpenID  string `json:"openID" bson:"openID"`
	UnionID string `json:"unionID" bson:"unionID"`
}

type SchoolInfo struct {
	Name   string `json:"name" bson:"name"`
	Status int    `json:"status" bson:"status"`
}

type RegionInfo struct {
	GaokaoRegion string `json:"gaokaoRegion" bson:"gaokaoRegion"`
}

type UserPublicMeta struct {
	School SchoolInfo `json:"school" bson:"school"`
}

type UserFreeToModifyMeta struct {
	Region RegionInfo `json:"region" bson:"region"`
}

type User struct {
	OID              primitive.ObjectID   `json:"oid" bson:"_id,omitempty"`
	Role             string               `json:"role" bson:"role"`
	Name             string               `json:"name" bson:"name"`
	EMail            string               `json:"email" bson:"email"`
	PasswordHash     string               `json:"-" bson:"passwordHash"`
	Avatar           string               `json:"avatar" bson:"avatar"`
	WechatInfo       WechatInfo           `json:"wechatInfo" bson:"wechatInfo"`
	PublicMeta       UserPublicMeta       `json:"publicMeta" bson:"publicMeta"`
	FreeToModifyMeta UserFreeToModifyMeta `json:"freeToModifyMeta" bson:"freeToModifyMeta"`
}

type UserPublicInfo struct {
	OID              string               `json:"id"`
	Avatar           string               `json:"avatar" bson:"avatar"`
	Name             string               `json:"name" bson:"name"`
	Role             string               `json:"role"`
	PublicMeta       UserPublicMeta       `json:"publicMeta"`
	FreeToModifyMeta UserFreeToModifyMeta `json:"freeToModifyMeta" bson:"freeToModifyMeta"`
}

func (user *User) ToPublicInfo() UserPublicInfo {
	return UserPublicInfo{
		OID:              user.OID.Hex(),
		Avatar:           user.Avatar,
		Name:             user.Name,
		Role:             user.Role,
		PublicMeta:       user.PublicMeta,
		FreeToModifyMeta: user.FreeToModifyMeta,
	}
}

type UserPublicInfoUpdateRequest struct {
	Avatar           string               `json:"avatar"`
	FreeToModifyMeta UserFreeToModifyMeta `json:"freeToModifyMeta"`
}

type LoginResponse struct {
	JWTToken string `json:"jwtToken"`
}

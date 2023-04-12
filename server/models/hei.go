package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// 学校办学模式
type HEIMode int

var PublicHEI HEIMode = 0                     // 公办机构
var PrivateHEI HEIMode = 1                    // 民办机构
var SinoForeignCooperative HEIMode = 2        // 中外合办
var Mainland_HK_MO_TW_Cooperative HEIMode = 3 // 大陆与港澳台合办

// 学校为专科或本科
type HEILevel int

var UniversityHEI HEILevel = 0 // 本科
var VocationalHEI HEILevel = 1 // 专科

// 学校地区
type Location struct {
	Provincial string `json:"provincial" bson:"provincial"` // 省级
	Municipal  string `json:"municipal" bson:"municipal"`   // 地级
}

// Higher education institute
type HEI struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code        string             `bson:"code" json:"code"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Location    Location           `bson:"location" json:"location"` // 地区
	Mode        HEIMode            `bson:"mode" json:"mode"`         // 办学模式
	Policies    []string           `bson:"policies" json:"policies"`
	Tags        []string           `bson:"tags" json:"tags"` // 学校称号
	Level       HEILevel           `bson:"level" json:"level"`
}

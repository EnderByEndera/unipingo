package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RoleAdmin      string = "ROLE_ADMIN"
	RoleUnpaidUser string = "ROLE_UNPAID_USER"
	RolePaidUser   string = "ROLE_PAID_USER"
)

type EducationalStageType string

// 注意，学位和学段都用以下变量
const (
	DegreeJunior   EducationalStageType = "JUNIOR_COLLEGE"
	DegreeBachelor EducationalStageType = "BACHELOR"
	DegreeMaster   EducationalStageType = "MASTER"
	DegreeDoctor   EducationalStageType = "DOCTOR"
	DegreePost     EducationalStageType = "POST"
)

type CollectionType string

const (
	CollectionItemMajor CollectionType = "COLLECTION_ITEM_MAJOR"
	CollectionItemHEI   CollectionType = "COLLECTION_ITEM_HEI"
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

type EduBGItem struct {
	AuthProcID     primitive.ObjectID   `json:"authProcID" bson:"authProcID"`
	MajorID        primitive.ObjectID   `json:"majorID" bson:"majorID"`
	MajorName      string               `json:"majorName" bson:"majorName"`
	HEIID          primitive.ObjectID   `json:"heiID" bson:"heiID"`
	HEIName        string               `json:"heiName" bson:"heiName"`
	Stage          EducationalStageType `json:"stage" bson:"stage"`
	EnrollmentTime uint64               `json:"enrollmentTime" bson:"enrollmentTime"`
	GraduationTime uint64               `json:"graduationTime" bson:"graduationTime"`
}

type CollectionItem struct {
	Type string             `json:"type" bson:"type"`
	ID   primitive.ObjectID `json:"id" bson:"id"`
	Text string             `json:"text" bson:"text"`
}

type AnswerInCollection struct {
	BelongsTo EntityWithName     `json:"answers" bson:"answers"`
	ID        primitive.ObjectID `json:"id" bson:"id"`
}

// 用户收藏的所有事项的结构体
type Collections struct {
	Majors  []EntityWithName     `json:"majors" bson:"majors"`
	HEIs    []EntityWithName     `json:"heis" bson:"heis"`
	Answers []AnswerInCollection `json:"answers" bson:"answers"`
}

// UserTag 用户标签
type UserTag struct {
	HEIName   string   `json:"heiName" bson:"heiName"`
	MajorName string   `json:"majorName" bson:"majorName"`
	CustomTag []string `json:"customTag" bson:"customTag"`
}

func (c *Collections) MarshalBSON() ([]byte, error) {
	if c.Answers == nil {
		c.Answers = make([]AnswerInCollection, 0)
	}
	if c.HEIs == nil {
		c.HEIs = make([]EntityWithName, 0)
	}
	if c.Majors == nil {
		c.Majors = make([]EntityWithName, 0)
	}
	type _t Collections
	return bson.Marshal((*_t)(c))
}

type User struct {
	ID                    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Role                  string             `json:"role" bson:"role"`
	Name                  string             `json:"name" bson:"name"`
	RealName              string             `json:"realName" bson:"realName"` // 真实姓名
	EMail                 string             `json:"email" bson:"email"`
	PasswordHash          string             `json:"-" bson:"passwordHash"`
	Avatar                string             `json:"avatar" bson:"avatar"`
	WechatInfo            WechatInfo         `json:"wechatInfo" bson:"wechatInfo"`
	Gender                string             `json:"gender" bson:"gender"`
	Birthday              uint64             `json:"birthday" bson:"birthday"`
	EducationalBackground []EduBGItem        `json:"educationalBackground" bson:"educationalBackground"`
	Collection            Collections        `json:"collection" bson:"collection"`
	Type                  string             `json:"type" bson:"type"`
	Membership            string             `json:"membership" bson:"membership"`
	UserTags              map[string]string  `json:"userTags" bson:"userTags"` // {"科目":"历史"}, {"爱好":"唱歌"}
	// PublicMeta   UserPublicMeta     `json:"publicMeta" bson:"publicMeta"`
	// FreeToModifyMeta UserFreeToModifyMeta `json:"freeToModifyMeta" bson:"freeToModifyMeta"`
}

type UserPublicInfo struct {
	OID              string               `json:"id"`
	Avatar           string               `json:"avatar" bson:"avatar"`
	Name             string               `json:"name" bson:"name"`
	Role             string               `json:"role"`
	PublicMeta       UserPublicMeta       `json:"publicMeta"`
	FreeToModifyMeta UserFreeToModifyMeta `json:"freeToModifyMeta" bson:"freeToModifyMeta"`
}

func (user *User) MarshalBSON() ([]byte, error) {
	if user.EducationalBackground == nil {
		user.EducationalBackground = make([]EduBGItem, 0)
	}
	type _u User
	return bson.Marshal((*_u)(user))
}

func (user *User) ToPublicInfo() UserPublicInfo {
	return UserPublicInfo{
		OID:    user.ID.Hex(),
		Avatar: user.Avatar,
		Name:   user.Name,
		Role:   user.Role,
		// PublicMeta: user.PublicMeta,
		// FreeToModifyMeta: user.FreeToModifyMeta,
	}
}

type UserPublicInfoUpdateRequest struct {
	Avatar string `json:"avatar" bson:"avatar"`
	Name   string `json:"name" bson:"name"`
	// FreeToModifyMeta UserFreeToModifyMeta `json:"freeToModifyMeta"`
}

type LoginResponse struct {
	JWTToken string `json:"jwtToken"`
}

type UserTagsInfoUpdateRequest struct {
	UserTags []string `json:"userTags"`
}

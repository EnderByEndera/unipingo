package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PersonalInfo 个人展示内容
type PersonalInfo struct {
	CEEPlace  string `bson:"CEEPlace" json:"CEEPlace"`             // 高考所在地
	Subject   string `bson:"subject" json:"subject"`               // 高考科目
	Age       uint32 `bson:"age" json:"age"`                       // 年龄
	Gender    string `bson:"gender" json:"gender"`                 // 性别
	Situation string `bson:"situation" json:"situation,omitempty"` // 具体情况
}

// QuestionBoxQuestion 提问箱问题
type QuestionBoxQuestion struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID   `bson:"userID" json:"userID"`           // 提问人ID
	Title       string               `bson:"title" json:"title"`             // 问题标题
	Description string               `bson:"description" json:"description"` // 问题描述
	CreateTime  uint64               `bson:"createTime" json:"createTime"`
	UpdateTime  uint64               `bson:"updateTime" json:"updateTime"`
	School      EntityWithName       `bson:"school" json:"school"`         // 所提问学校的ID
	Major       EntityWithName       `bson:"major" json:"major"`           // 所提问专业的ID
	Questioner  PersonalInfo         `bson:"questioner" json:"questioner"` // 提问者相关数据
	Answers     []primitive.ObjectID `bson:"answers" json:"answers"`       // 该问题下所有的回答
	AskTo       []primitive.ObjectID `bson:"askTo" json:"askTo"`           // 该问题向谁提问
	AskTags     []string             `bson:"askTags" json:"askTags"`       // 该问题向带有哪些标签的用户提问
}

func (question *QuestionBoxQuestion) Init() {
	question.CreateTime = uint64(time.Now().Unix())
	question.Answers = make([]primitive.ObjectID, 0)
}

type QuestionInLabelInfo struct {
	ID   primitive.ObjectID `bson:"questionID" json:"questionID"`
	Name string             `bson:"name" json:"name"`
}

// QuestionBoxLabel 提问箱标签
type QuestionBoxLabel struct {
	ID         primitive.ObjectID    `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID    `bson:"userID" json:"userID"`
	Content    string                `bson:"content" json:"content"`
	CreateTime uint64                `bson:"createTime" json:"createTime"`
	UpdateTime uint64                `bson:"updateTime" json:"updateTime"`
	Questions  []QuestionInLabelInfo `bson:"questions" json:"questions"`
}

func (ql *QuestionBoxLabel) Init() {
	if ql.CreateTime == 0 {
		ql.CreateTime = uint64(time.Now().Unix())
		ql.UpdateTime = uint64(time.Now().Unix())
	}
}

// QuestionBoxAnswer 提问箱回答
type QuestionBoxAnswer struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"userID" json:"userID"`         // 回答用户ID
	Content    string             `bson:"content" json:"content"`       // 回答内容
	QuestionID primitive.ObjectID `bson:"questionID" json:"questionID"` // 所属问题的ID
	School     EntityWithName     `bson:"school" json:"school"`
	Major      EntityWithName     `bson:"major" json:"major"`
	Statistics AnswerStats        `bson:"answerStats" json:"answerStats"` // 回答相关的数据
	CreateTime uint64             `bson:"createTime" json:"createTime"`
	UpdateTime uint64             `bson:"updateTime" json:"updateTime"`
	Respondant PersonalInfo       `bson:"respondant" json:"respondant"` // 回答者相关数据
}

type QuestionBoxAnswerReq struct {
	Content    string             `json:"content"`    // 回答内容
	QuestionID primitive.ObjectID `json:"questionID"` // 所属问题的ID
	School     EntityWithName     `json:"school"`
	Major      EntityWithName     `json:"major"`
	Respondant PersonalInfo       `json:"respondant"` // 回答者相关数据
}

func (answer *QuestionBoxAnswer) Init() {
	answer.CreateTime = uint64(time.Now().Unix())
	answer.UpdateTime = uint64(time.Now().Unix())
}

type AnswerLog struct {
	AnswerID   primitive.ObjectID `bson:"answerID" json:"answerID"`
	QuestionID primitive.ObjectID `bson:"questionID" json:"questionID"`
	LogTime    uint64             `bson:"logTime" json:"logTime"`
}

type AnswerLogs struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Answers []AnswerLog        `bson:"answers" json:"answers"`
}

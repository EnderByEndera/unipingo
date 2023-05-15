package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 个人展示内容
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
}

func (question *QuestionBoxQuestion) Init() {
	question.CreateTime = uint64(time.Now().Unix())
	question.Answers = make([]primitive.ObjectID, 0)
}

type NewQuestionReq struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	School      EntityWithName `json:"school"`
	Major       EntityWithName `json:"major"`
	Questioner  PersonalInfo   `json:"questioner"`
}

type QuestionInLabelInfo struct {
	ID      primitive.ObjectID `bson:"questionID" json:"questionID"`
	Name    string             `bson:"name" json:"name"`
	HasRead bool               `bson:"hasRead" json:"hasRead"`
}

type LabelStats struct {
	ReadNum     uint32 `bson:"readNum" json:"readNum"`
	QuestionNum uint32 `bson:"questionNum" json:"questionNum"`
}

// 提问箱标签
type QuestionLabel struct {
	ID         primitive.ObjectID    `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID    `bson:"userID" json:"userID"`
	Content    string                `bson:"content" json:"content"`
	CreateTime uint64                `bson:"createTime" json:"createTime"`
	UpdateTime uint64                `bson:"updateTime" json:"updateTime"`
	Statistics LabelStats            `bson:"stats" json:"stats"` // 问题标签相关的数据
	Questions  []QuestionInLabelInfo `bson:"questions" json:"questions"`
}

func (ql *QuestionLabel) Init() {
	ql.CreateTime = uint64(time.Now().Unix())
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
}

type GetLabelsFromUserRequest struct {
	Page    int64 `json:"page"`
	PageNum int64 `json:"pageNum"`
}

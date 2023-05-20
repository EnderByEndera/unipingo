package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"melodie-site/server/models"
)

type NewQuestionReq struct {
	Title       string                `json:"title"`       // 问题标题
	Description string                `json:"description"` // 问题描述
	School      models.EntityWithName `json:"school"`      // 提问学校
	Major       models.EntityWithName `json:"major"`       // 提问专业
	Questioner  models.PersonalInfo   `json:"questioner"`  // 提问人信息
}

type NewQuestionRes struct {
	QuestionID primitive.ObjectID `json:"questionID"` // 新增问题ID
}

type QueryQuestionByIDRes struct {
	Question *models.QuestionBoxQuestion `json:"question"`
}

type QueryQuestionListRes struct {
	Questions []*models.QuestionBoxQuestion `json:"questions"`
	NextPage  int                           `json:"next_page"`
}

type UpdateQuestionDescriptionReq struct {
	Question *models.QuestionBoxQuestion `json:"question"`
}

type UpdateQuestionDescriptionRes struct {
	Update bool `json:"update"`
}

type UpdateQuestionSchoolOrMajorReq struct {
	Question *models.QuestionBoxQuestion `json:"question"`
}

type UpdateQuestionSchoolOrMajorRes struct {
	Update bool `json:"update"`
}

type NewLabelsReq struct {
	Labels []*models.QuestionLabel `json:"labels"`
}

type NewLabelsRes struct {
	LabelIDs []primitive.ObjectID `json:"labelIDs"`
}

type GetLabelsFromUserReq struct {
}

type GetLabelsFromUserRes struct {
	Labels []*models.QuestionLabel `json:"labels"`
}

type GetLabelsFromQuestionReq struct {
	QuestionID primitive.ObjectID `json:"questionID"`
}

type GetLabelsFromQuestionRes struct {
	Labels []*models.QuestionLabel `json:"labels"`
}

type DeleteLabelReq struct {
}

type DeleteLabelRes struct {
	LabelID primitive.ObjectID `json:"labelID"`
}

type UpdateLabelContentReq struct {
	LabelID primitive.ObjectID `json:"labelID"`
	Content string             `json:"content"`
}

type UpdateLabelContentRes struct {
	LabelID primitive.ObjectID `json:"labelID"`
}

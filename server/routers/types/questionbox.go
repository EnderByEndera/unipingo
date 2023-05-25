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
	AskTo       []primitive.ObjectID  `json:"askTo"`       // 该问题向谁提问
	AskTags     []string              `json:"astTags"`     // 该问题向带有哪些标签的用户提问
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
	Labels   []*models.QuestionBoxLabel  `json:"labels"`
	Question *models.QuestionBoxQuestion `json:"question"`
}

type NewLabelsRes struct {
	LabelIDs []primitive.ObjectID `json:"labelIDs"`
}

type GetLabelsFromUserReq struct {
}

type GetLabelsFromUserRes struct {
	Labels []*models.QuestionBoxLabel `json:"labels"`
}

type GetLabelsFromQuestionReq struct {
	QuestionID primitive.ObjectID `json:"questionID"`
}

type GetLabelsFromQuestionRes struct {
	Labels []*models.QuestionBoxLabel `json:"labels"`
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

type ReadAnswerByUserReq struct {
	AnswerID primitive.ObjectID `json:"answerID"`
}

type ReadAnswerByUserRes struct {
	Log bool `json:"log"`
}

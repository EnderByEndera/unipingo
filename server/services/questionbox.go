package services

import (
	"melodie-site/server/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionBoxService struct {
}

var (
	questionBoxService *QuestionBoxService
)

func GetQuestionBoxService() *QuestionBoxService {
	if questionBoxService != nil {
		return questionBoxService
	}

	questionBoxService = &QuestionBoxService{}
	return questionBoxService
}

func (service *QuestionBoxService) NewQuestion(question *models.QuestionBoxQuestion, label *models.QuestionLabel) (err error) {
	// TODO: 需要根据接受的问题和标签创建新问题
	return
}

func (service *QuestionBoxService) UpdateQuestion(question *models.QuestionBoxQuestion) (updatedQuestion *models.QuestionBoxQuestion, err error) {
	return
}

func (service *QuestionBoxService) QueryQuestionByID(questionID primitive.ObjectID) (question *models.QuestionBoxQuestion, err error) {
	return
}

func (service *QuestionBoxService) QuestionList(user *models.User, page uint64) (questions []*models.QuestionBoxQuestion, err error) {
	return
}

func (service *QuestionBoxService) NewAnswer(answer *models.Answer) (err error) {
	return
}

func (service *QuestionBoxService) AnswerList(question *models.QuestionBoxQuestion, page uint64) (answers []*models.QuestionBoxAnswer, err error) {
	return
}

package services

import (
	"context"
	"melodie-site/server/db"
	"melodie-site/server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (service *QuestionBoxService) NewAnswer(answer *models.QuestionBoxAnswer) (err error) {
	answer.Init()
	//TODO:数据库名字
	_, err = db.GetCollection("questionboxanswer").InsertOne(context.Background(), answer)
	if err != nil {
		return
	}
	return
}

func (service *QuestionBoxService) QueryAnswerByID(answerID primitive.ObjectID) (answer *models.QuestionBoxAnswer, err error) {
	filter := bson.M{"_id": answerID}
	answer = &models.QuestionBoxAnswer{}
	err = db.GetCollection("questionboxanswer").FindOne(context.TODO(), filter).Decode(answer)
	return
}

func (service *AnswerService) DeleteQuestionboxAnswerByID(answerID primitive.ObjectID) (err error) {
	filter := bson.D{{"_id", answerID}}
	_, err = db.GetCollection("questionboxanswer").DeleteOne(context.TODO(), filter)
	return
}

func (service *QuestionBoxService) AnswerList(question *models.QuestionBoxQuestion, page int64) (answers []*models.QuestionBoxAnswer, err error) {
	questionID := question.ID
	filter := bson.M{"questionID": questionID}
	opts := options.Find().SetLimit(20).SetSkip(20 * page)
	// TODO: 数据库名字
	res, err := db.GetCollection("questionboxanswer").Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}
	err = res.All(context.TODO(), &answers)
	return
}

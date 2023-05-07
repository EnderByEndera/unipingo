package services

import (
	"context"
	"errors"
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

func questionExists(question *models.QuestionBoxQuestion) bool {
	// 只要在相同学校或专业下存在相同问题则判定为True
	return db.GetCollection("questions").FindOne(context.TODO(), bson.M{
		"title":  question.Title,
		"school": question.School,
		"major":  question.Major,
	}).Err() == nil
}

func (service *QuestionBoxService) NewQuestion(question *models.QuestionBoxQuestion) (docID primitive.ObjectID, err error) {
	if question.Title == "" || question.Description == "" {
		err = errors.New("该问题没有填写标题或描述")
		return
	}
	if questionExists(question) {
		err = errors.New("该问题已存在")
		return
	}
	question.Init()
	res, err := db.GetCollection("questions").InsertOne(context.TODO(), question)
	if err != nil {
		return
	}
	docID = res.InsertedID.(primitive.ObjectID)
	return
}

func (service *QuestionBoxService) UpdateQuestionDescription(question *models.QuestionBoxQuestion) (err error) {
	if question.Description == "" {
		err = errors.New("更新描述为空")
		return
	}
	filter := bson.M{
		"_id": question.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"description": question.Description,
		},
	}
	err = db.GetCollection("questions").FindOneAndUpdate(context.TODO(), filter, update).Err()
	return
}

func (service *QuestionBoxService) UpdateQuestionSchoolOrMajor(question *models.QuestionBoxQuestion) (err error) {
	if question.School.ID == primitive.NilObjectID && question.Major.ID == primitive.NilObjectID {
		err = errors.New("更新学校或专业均为空")
		return
	}

	filter := bson.M{
		"_id": question.ID,
	}
	update := bson.M{
		"$set": bson.M{
			"school": question.School,
			"major":  question.Major,
		},
	}
	err = db.GetCollection("questions").FindOneAndUpdate(context.TODO(), filter, update).Err()
	return
}

func (service *QuestionBoxService) QueryQuestionByID(questionID primitive.ObjectID) (question *models.QuestionBoxQuestion, err error) {
	question = new(models.QuestionBoxQuestion)
	if questionID == primitive.NilObjectID {
		err = errors.New("questionID为空")
		return
	}
	err = db.GetCollection("questions").FindOne(context.TODO(), bson.M{"_id": questionID}).Decode(question)
	return
}

func (service *QuestionBoxService) QuestionList(user *models.User, page int64, pageNum int64) (questions []*models.QuestionBoxQuestion, err error) {
	if user == nil {
		err = errors.New("user为空")
		return
	}

	filter := bson.M{
		"userID": user.ID,
	}

	if page < 0 || pageNum < 0 {
		err = errors.New("page或pageNum小于0")
		return
	}
	opts := options.Find().SetLimit(pageNum).SetSkip(page * pageNum)
	cur, err := db.GetCollection("questions").Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}

	cur.All(context.TODO(), &questions)
	return
}

func (service *QuestionBoxService) AddAnswerToQuestion(questionID primitive.ObjectID, answer *models.QuestionBoxAnswer) (err error) {
	filter := bson.M{
		"_id": questionID,
	}
	update := bson.M{
		"$push": bson.M{
			"answers": answer,
		},
	}
	err = db.GetCollection("questions").FindOneAndUpdate(context.TODO(), filter, update).Err()
	return
}

func (service *QuestionBoxService) NewLabels(labels []*models.QuestionLabel) (labelIDs []primitive.ObjectID, err error) {
	// 如果问题不存在标签，则直接退出
	if labels == nil {
		return
	}

	for _, label := range labels {
		if label.Content == "" {
			err = errors.New("部分标签没有描述信息")
			return
		}

		filter := bson.M{
			"userID":  label.UserID,
			"content": label.Content,
		}
		update := bson.D{
			{Key: "$push", Value: bson.D{{"questions", label.Questions[0]}}}, // 如果该数据库中存在该标签，则将该标签关联的问题添加到数据库中
			{Key: "$inc", Value: bson.D{{"stats.questionNum", 1}}},           // 如果该数据库中存在该标签，则将该标签对应的问题数量+1
		}
		opts := options.FindOneAndUpdate().SetUpsert(true) // 如果不存在该标签，则将该标签添加到数据库中
		err = db.GetCollection("labels").FindOneAndUpdate(context.TODO(), filter, update, opts).Err()
		if err != nil {
			return
		}
	}
	return
}

func (service *QuestionBoxService) QueryQuestionsFromLabelID(labelID primitive.ObjectID, page, pageNum int64) (questions []*models.QuestionBoxQuestion, err error) {
	filter := bson.M{
		"_id": labelID,
	}
	label := &models.QuestionLabel{}
	err = db.GetCollection("labels").FindOne(context.TODO(), filter).Decode(label)
	if err != nil {
		return
	}
	questionInfos := label.Questions[page*pageNum : (page+1)*pageNum]
	questionIDs := make([]primitive.ObjectID, 0)
	for _, info := range questionInfos {
		questionIDs = append(questionIDs, info.ID)
	}
	qInfoFilter := bson.M{
		"_id": bson.M{
			"$in": questionIDs,
		},
	}
	cur, err := db.GetCollection("questions").Find(context.TODO(), qInfoFilter)
	if err != nil {
		return
	}

	cur.All(context.TODO(), &questions)
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

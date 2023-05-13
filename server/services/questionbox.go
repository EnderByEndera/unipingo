package services

import (
	"context"
	"errors"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"time"

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

	err = cur.All(context.TODO(), &questions)
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

	err = cur.All(context.TODO(), &questions)
	return
}

func (service *QuestionBoxService) AddAnswerToQuestion(questionID primitive.ObjectID, answer *models.QuestionBoxAnswer) (err error) {
	filter := bson.M{
		"_id": questionID,
	}
	update := bson.M{
		"$addToSet": bson.M{
			"answers": answer,
		},
	}
	err = db.GetCollection("questions").FindOneAndUpdate(context.TODO(), filter, update).Err()
	return
}

func (service *QuestionBoxService) DeleteAnswerFromQuestion(questionID primitive.ObjectID, answerID primitive.ObjectID) (err error) {
	return
}

func (service *QuestionBoxService) DeleteQuestion(questionID primitive.ObjectID) (err error) {
	err = db.GetCollection("questions").FindOneAndDelete(context.TODO(), bson.M{"_id": questionID}).Err()
	if err != nil {
		return
	}

	update := bson.M{
		"$inc": bson.M{
			"stats.questionNum": 1,
		},
		"$pull": bson.M{
			"questions.questionID": questionID,
		},
	}
	_, err = db.GetCollection("labels").UpdateMany(context.TODO(), bson.M{"questions.questionID": questionID}, update)
	return
}

func (service *QuestionBoxService) NewLabels(labels []*models.QuestionLabel) (newLabels []models.EntityWithName, err error) {
	// 如果问题不存在标签，则直接退出
	if labels == nil {
		return
	}

	for _, label := range labels {
		if label.Content == "" {
			err = errors.New("部分标签没有描述信息")
			return
		}

		if label.UserID == primitive.NilObjectID {
			err = errors.New("部分标签不存在用户信息")
			return
		}

		filter := bson.M{
			"userID":  label.UserID,
			"content": label.Content,
		}

		err = db.GetCollection("labels").FindOne(context.TODO(), filter).Err()
		if err != nil {
			label.Init()
			res, labelErr := db.GetCollection("labels").InsertOne(context.TODO(), label)
			if labelErr != nil {
				err = labelErr
				return
			}
			if labelID, ok := res.InsertedID.(primitive.ObjectID); ok {
				newLabels = append(newLabels, models.EntityWithName{ID: labelID, Name: label.Content})
			}
		}
		err = nil
	}
	return
}

func (service *QuestionBoxService) QueryLabelByID(labelID primitive.ObjectID) (label *models.QuestionLabel, err error) {
	return
}

func (service *QuestionBoxService) QueryLabelsFromUser(user *models.User, page, pageNum int64) (labels []*models.QuestionLabel, err error) {
	filter := bson.M{
		"userID": user.ID,
	}

	opts := options.Find().SetLimit(pageNum).SetSkip(page * pageNum)

	cur, err := db.GetCollection("labels").Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}

	err = cur.All(db.GetMongoConn().Context, &labels)
	return
}

func (service *QuestionBoxService) QueryLabelFromQuestion(user *models.User, question *models.QuestionBoxQuestion, page, pageNum int64) (labels []*models.QuestionLabel, err error) {
	filter := bson.M{
		"userID": user.ID,
		"questions": bson.M{
			"$elemMatch": bson.M{
				"questionID": question.ID,
			},
		},
	}

	opts := options.Find().SetLimit(pageNum).SetSkip(page * pageNum)
	cur, err := db.GetCollection("labels").Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}

	err = cur.All(context.TODO(), &labels)
	return
}

func (service *QuestionBoxService) UpdateLabelContent(label *models.QuestionLabel) (err error) {
	err = db.GetCollection("labels").FindOneAndUpdate(context.TODO(),
		bson.M{"_id": label.ID},
		bson.M{
			"$set": bson.M{
				"content":    label.Content,
				"updateTime": uint64(time.Now().Unix()),
			}}).Err()
	return
}

func (service *QuestionBoxService) AddQuestionInLabel(labelID primitive.ObjectID, question *models.QuestionInLabelInfo) (err error) {
	filter := bson.M{
		"_id": labelID,
	}

	update := bson.M{
		"$inc": bson.M{
			"stats.questionNum": 1,
		},
		"$addToSet": bson.M{
			"questions": question,
		},
	}

	err = db.GetCollection("labels").FindOneAndUpdate(context.TODO(), filter, update).Err()
	return
}

func (service *QuestionBoxService) DeleteQuestionFromLabel(labelID primitive.ObjectID, questionID primitive.ObjectID) (err error) {
	filter := bson.M{
		"labelID": labelID,
	}

	update := bson.M{
		"$pull": bson.M{
			"questions.questionID": questionID,
		},
		"$inc": bson.M{
			"stats.questionNum": 1,
		},
	}

	err = db.GetCollection("labels").FindOneAndUpdate(context.TODO(), filter, update).Err()
	return
}

func (service *QuestionBoxService) DeleteLabel(label *models.QuestionLabel) (err error) {
	err = db.GetCollection("labels").FindOneAndDelete(context.TODO(), bson.M{"_id": label.ID}).Err()
	return
}

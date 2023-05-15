package services

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"time"

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

/*
--------------------------------------------
问题模块区域
*/

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
	// TODO 关于questionexits是否要存在
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
	questionIDs := make([]primitive.ObjectID, 0)

	transaction := func(sessCtx mongo.SessionContext) (res interface{}, tranErr error) {
		tranErr = db.GetCollection("labels").FindOne(sessCtx, filter).Decode(label)

		if tranErr != nil {
			return
		}
		questionInfos := label.Questions[page*pageNum : (page+1)*pageNum]

		for _, info := range questionInfos {
			questionIDs = append(questionIDs, info.ID)
		}
		qInfoFilter := bson.M{
			"_id": bson.M{
				"$in": questionIDs,
			},
		}

		cur, tranErr := db.GetCollection("questions").Find(sessCtx, qInfoFilter)
		if tranErr != nil {
			return
		}

		tranErr = cur.All(sessCtx, &questions)
		return
	}
	_, err = db.GetMongoConn().UseSession(nil, transaction)
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

func (service *QuestionBoxService) AddAnswerToQuestion(questionID primitive.ObjectID, answerID primitive.ObjectID) (err error) {
	filter := bson.M{
		"_id": questionID,
	}

	update := bson.M{
		"$addToSet": bson.M{
			"answers": answerID,
		},
		"$set": bson.M{
			"updateTime": uint64(time.Now().Unix()),
		},
	}

	err = db.GetCollection("questions").FindOneAndUpdate(context.TODO(), filter, update).Err()
	return
}

func (service *QuestionBoxService) DeleteQuestion(questionID primitive.ObjectID) (err error) {
	transaction := func(sessCtx mongo.SessionContext) (res interface{}, err error) {
		question := new(models.QuestionBoxQuestion)
		err = db.GetCollection("questions").FindOneAndDelete(context.TODO(), bson.M{"_id": questionID}).Decode(question)
		if err != nil {
			return
		}

		update := bson.M{
			"$inc": bson.M{
				"stats.questionNum": -1,
			},
			"$pull": bson.M{
				"questions.questionID": questionID,
			},
		}
		_, err = db.GetCollection("labels").UpdateMany(context.TODO(), bson.M{"questions.questionID": questionID}, update)
		return
	}
	_, sessErr := db.GetMongoConn().UseSession(nil, transaction)
	return sessErr
}

/*
--------------------------------------------------
标签（文件夹）模块区域
*/

// NewLabels 创建多个标签
func (service *QuestionBoxService) NewLabels(labels []*models.QuestionLabel) (newLabels []models.EntityWithName, err error) {
	// 如果问题不存在标签，则直接退出
	if labels == nil {
		return
	}

	for _, label := range labels {
		if label.Content == "" {
			err = fmt.Errorf("部分标签不存在内容")
			return
		}

		if label.UserID == primitive.NilObjectID {
			err = fmt.Errorf("标签%s不存在用户信息", label.Content)
			return
		}

		filter := bson.M{
			"userID":  label.UserID,
			"content": label.Content,
		}

		_, err = db.GetMongoConn().UseSession(nil, func(sessCtx mongo.SessionContext) (interface{}, error) {
			findErr := db.GetCollection("labels").FindOne(sessCtx, filter).Err()
			if findErr != nil {
				label.Init()
				res, labelErr := db.GetCollection("labels").InsertOne(sessCtx, label)
				if labelErr != nil {
					return nil, labelErr
				}
				if labelID, ok := res.InsertedID.(primitive.ObjectID); ok {
					newLabels = append(newLabels, models.EntityWithName{ID: labelID, Name: label.Content})
				}
			}

			return nil, nil
		})

		if err != nil {
			return
		}
	}
	return
}

func (service *QuestionBoxService) QueryLabelByID(labelID primitive.ObjectID) (label *models.QuestionLabel, err error) {
	label = new(models.QuestionLabel)
	err = db.GetCollection("labels").FindOne(context.TODO(), bson.M{"_id": labelID}).Decode(label)
	return
}

func (service *QuestionBoxService) QueryLabelsFromUser(user *models.User, page, pageNum int64) (labels []*models.QuestionLabel, err error) {
	filter := bson.M{
		"userID": user.ID,
	}

	opts := options.Find().SetLimit(pageNum).SetSkip(page * pageNum)

	_, err = db.GetMongoConn().UseSession(nil, func(sessCtx mongo.SessionContext) (interface{}, error) {
		cur, err := db.GetCollection("labels").Find(sessCtx, filter, opts)
		if err != nil {
			return nil, err
		}

		err = cur.All(sessCtx, &labels)
		return nil, err
	})
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

	_, err = db.GetMongoConn().UseSession(nil, func(ctx mongo.SessionContext) (res interface{}, err error) {
		opts := options.Find().SetLimit(pageNum).SetSkip(page * pageNum)
		cur, err := db.GetCollection("labels").Find(ctx, filter, opts)
		if err != nil {
			return
		}

		err = cur.All(ctx, &labels)
		return
	})
	return
}

func (service *QuestionBoxService) DeleteLabel(label *models.QuestionLabel) (err error) {
	err = db.GetCollection("labels").FindOneAndDelete(context.TODO(), bson.M{"_id": label.ID}).Err()
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

	type Result struct {
		ID  primitive.ObjectID `bson:"_id,omitempty"`
		Cnt int                `bson:"cnt"`
	}

	_, err = db.GetMongoConn().UseSession(nil, func(sessCtx mongo.SessionContext) (res interface{}, tranErr error) {
		update := bson.M{
			"$addToSet": bson.M{
				"questions": question,
			},
			"$set": bson.M{
				"updateTime": uint64(time.Now().Unix()),
			},
		}

		tranErr = db.GetCollection("labels").FindOneAndUpdate(sessCtx, filter, update).Err()
		if tranErr != nil {
			return
		}

		cur, tranErr := db.GetCollection("labels").Aggregate(sessCtx, mongo.Pipeline{
			{{"$match", bson.M{"_id": labelID}}},
			{{"$project", bson.M{"cnt": bson.M{"$size": "$questions"}}}},
		})
		if tranErr != nil {
			return
		}

		result := make([]Result, 0)
		tranErr = cur.All(sessCtx, &result)
		if tranErr != nil {
			return
		}

		if len(result) == 0 {
			result = append(result, Result{Cnt: 0})
		}

		tranErr = db.GetCollection("labels").FindOneAndUpdate(sessCtx, filter, bson.M{
			"$set": bson.M{
				"stats.questionNum": result[0].Cnt,
			},
		}).Err()

		return
	})
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
		"$set": bson.M{
			"updateTime": uint64(time.Now().Unix()),
		},
	}

	err = db.GetCollection("labels").FindOneAndUpdate(context.TODO(), filter, update).Err()
	return
}

/*
--------------------------------------------------------
回答模块区域
*/

// NewAnswer 创建新回答
func (service *QuestionBoxService) NewAnswer(answer *models.QuestionBoxAnswer) (docID primitive.ObjectID, err error) {
	if answer.Content == "" {
		err = errors.New("该回答没有填写内容")
		return
	}
	answer.Init()

	questionID := answer.QuestionID
	//验证question存在
	_, err = questionBoxService.QueryQuestionByID(questionID)
	if err != nil {
		return
	}
	// TODO 改为原子操作
	res, err := db.GetCollection("questionboxanswer").InsertOne(context.Background(), answer)
	if err != nil {
		return
	}
	docID = res.InsertedID.(primitive.ObjectID)
	err = questionBoxService.AddAnswerToQuestion(questionID, docID)
	if err != nil {
		err = errors.New("回答和问题关联失败")
		return
	}

	return
}

func (service *QuestionBoxService) QueryAnswerByID(answerID primitive.ObjectID) (answer *models.QuestionBoxAnswer, err error) {
	filter := bson.M{"_id": answerID}
	answer = &models.QuestionBoxAnswer{}
	if answerID == primitive.NilObjectID {
		err = errors.New("answerID为空")
		return
	}
	err = db.GetCollection("questionboxanswer").FindOne(context.TODO(), filter).Decode(answer)
	return
}

func (service *QuestionBoxService) DeleteQuestionBoxAnswerByID(answerID primitive.ObjectID) (err error) {
	if answerID == primitive.NilObjectID {
		err = errors.New("answerID为空")
		return
	}
	filter := bson.M{"_id": answerID}
	_, err = db.GetCollection("questionboxanswer").DeleteOne(context.TODO(), filter)

	// TODO: 这里还需要将每个question里面的这个回答ID给删除掉才可以，MongoDB不支持外键就是这么麻烦 :-)
	return
}

// AnswerList 获取一个问题对应的所有回答
func (service *QuestionBoxService) AnswerList(question *models.QuestionBoxQuestion, page int64, pageNum int64) (answers []*models.QuestionBoxAnswer, err error) {
	if question == nil {
		err = errors.New("question为空")
		return
	}
	questionID := question.ID
	filter := bson.M{"questionID": questionID}

	if page < 0 || pageNum < 0 {
		err = errors.New("page或pageNum小于0")
		return
	}
	opts := options.Find().SetLimit(pageNum).SetSkip(pageNum * page)

	res, err := db.GetCollection("questionboxanswer").Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}
	err = res.All(context.TODO(), &answers)
	return
}

// MyAnswerList 获取当前用户的所有回答（提问箱部分的“我的回答”）
func (service *QuestionBoxService) MyAnswerList(user *models.User, page int64, pageNum int64) (answers []*models.QuestionBoxAnswer, err error) {
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
	cur, err := db.GetCollection("questionboxanswer").Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}

	err = cur.All(context.TODO(), &answers)
	return
}

func (service *QuestionBoxService) UpdateAnswerContent(answer *models.QuestionBoxAnswer) (err error) {
	if answer.Content == "" {
		err = errors.New("更新的回答为空")
		return
	}

	filter := bson.M{
		"_id": answer.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"content": answer.Content,
		},
	}
	err = db.GetCollection("questionboxanswer").FindOneAndUpdate(context.TODO(), filter, update).Err()
	return
}

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
	/* TODO :
	1、个人感觉这里是不是不用判定问题是否相同，因为提问箱应该是向具体的人或者特定群体提问的。
	即使问题一样，提问者可能选择的提问对象也会不同。
	其次，产品那边应该还暂时没有给出“如何看到其他人提问的问题”的设计。
	师兄您的意思可能是 如果问题已经存在就不能提问了，用户需要搜索看别人问过的问题和答案
	再者，如果后续提问是需要付费的，感觉每次提问都创建一个question，这样会方便order那边保存唯一的questionID
	2、questionExists这个函数中：（school、major、title都一样就确认问题已存在），
	个人感觉对同一个title所提出的具体问题也可能不一样的（即content不同），
	那这样是不是也不能算作是一个question呢？
	（我这里是将title理解为主题，就是类似“学习环境”，“宿舍环境”之类的）
	*/
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
	// TODO  我看见现在的figma上又更新了提问表单的内容，上边显示在提问的时候可以选择文件夹，所以感觉可能还要加点什么
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
	// TODO 需要验证问题的school、major和答案的school、major相同吗（如果answer的两个属性来自于question就不用了？）
	filter := bson.M{
		"_id": questionID,
	}

	// TODO 这个好像把整个answer给赋值了，我看model里面定义questions时answers里面存储的是所有问题的ID
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

// TODO 对具体内容没有问题，但是不太清楚这个函数的功能，是创建问题文件夹吗？然后labelIDs是在哪里返回的呢？
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
	/*TODO：这里是将answer和question关联，调用了一下师兄的AddAnswerToQuestion函数。

	这里InsertOne和AddAnswerToQuestion感觉类似一个事务，一个执行失败另一个也要回滚；比如
	这里如果已经执行InsertOne成功，但AddAnswerToQuestion执行失败，那这个回答即使存入数据库也是
	没有用的数据了（没有和question绑定，后续用questionID就查不出来）。

	如果这种没用的数据没有什么影响 或者这两个函数执行的成功率都很高以至于没用的数据很少，
	或许就这样写着？只让前端提示用户“提交回答失败”就行？
	*/
	res, err := db.GetCollection("questionboxanswer").InsertOne(context.Background(), answer)
	if err != nil {
		return
	}
	docID = res.InsertedID.(primitive.ObjectID)
	err = questionBoxService.AddAnswerToQuestion(questionID, answer)
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


func (service *AnswerService) DeleteQuestionboxAnswerByID(answerID primitive.ObjectID) (err error) {
	if answerID == primitive.NilObjectID {
		err = errors.New("answerID为空")
		return
	}
	filter := bson.M{"_id": answerID}
	_, err = db.GetCollection("questionboxanswer").DeleteOne(context.TODO(), filter)
	return
}

// 获取一个问题对应的所有回答
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

// 获取当前用户的所有回答（提问箱部分的“我的回答”）
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

	cur.All(context.TODO(), &answers)
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
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
	//如果已经存在  用户要到哪里去看呢（大学和专业那里展示的问题应该都是写死的，和提问箱无关）提问箱的问题只有参与者知道，其他用户都看不到
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
	// TODO 需要验证问题的school、major和答案的school、major相同吗（如果answer的两个属性来自于question就不用了吧）
	filter := bson.M{
		"_id": questionID,
	}
	//这个好像把answer所有属性值给answers属性了
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

func (service *QuestionBoxService) NewAnswer(answer *models.QuestionBoxAnswer) (docID primitive.ObjectID, err error) {
	if answer.Content == "" {
		err = errors.New("该回答没有填写内容")
		return
	}
	answer.Init()

	questionID := answer.QuestionID
	//和question关联
	//TODO：如果后面insertOne插入错误，AddAnswerToQuestion函数需要回滚
	_, err = questionBoxService.QueryQuestionByID(questionID)
	if err != nil {
		return
	}
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

// 这个函数是列出一个问题的所有回答
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

// 这个函数列出我的回答列表
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

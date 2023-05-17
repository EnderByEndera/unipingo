package services

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type QuestionBoxService struct {
	wc *writeconcern.WriteConcern
	rc *readconcern.ReadConcern
}

var (
	questionBoxService *QuestionBoxService
)

func GetQuestionBoxService() *QuestionBoxService {
	if questionBoxService != nil {
		return questionBoxService
	}

	questionBoxService = &QuestionBoxService{
		// Set w = 1 to ensure the maximum performance
		wc: writeconcern.New(writeconcern.W(1), writeconcern.J(true)),
		rc: readconcern.New(readconcern.Level("local")),
	}
	return questionBoxService
}

/*
--------------------------------------------
问题模块区域
*/

func questionExists(question *models.QuestionBoxQuestion) (ok bool) {
	// 只要在相同学校或专业下存在相同问题则判定为True
	ok = db.GetCollection("questions").FindOne(context.TODO(), bson.M{
		"title":  question.Title,
		"userID": question.UserID,
	}).Decode(question) == nil
	return
}

func (service *QuestionBoxService) NewQuestion(question *models.QuestionBoxQuestion) (questionID primitive.ObjectID, err error) {
	if question.Title == "" || question.Description == "" {
		err = errors.New("该问题没有填写标题或描述")
		return
	}

	if question.School.ID == primitive.NilObjectID && question.Major.ID == primitive.NilObjectID {
		err = errors.New("该问题学校和专业均为空")
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
		questionID = question.ID
		return
	}
	question.Init()
	res, err := db.GetCollection("questions").InsertOne(context.TODO(), question)
	if err != nil {
		return
	}

	questionID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		err = errors.New("类型断言失败")
	}
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
	sessOpts := options.Session().SetDefaultWriteConcern(service.wc).SetDefaultReadConcern(service.rc)
	_, err = db.GetMongoConn().UseSession(sessOpts, transaction)
	return
}

func (service *QuestionBoxService) QueryQuestionsFromUser(user *models.User, page int64, pageNum int64) (questions []*models.QuestionBoxQuestion, sessErr error) {
	if user == nil {
		sessErr = errors.New("user为空")
		return
	}

	filter := bson.M{
		"userID": user.ID,
	}

	if page < 0 || pageNum <= 0 {
		sessErr = fmt.Errorf("page: %d, pageNum: %d", page, pageNum)
		return
	}
	opts := options.Find().SetLimit(pageNum).SetSkip(page * pageNum)

	_, sessErr = db.GetMongoConn().UseSession(nil, func(sessCtx mongo.SessionContext) (interface{}, error) {
		cur, err := db.GetCollection("questions").Find(sessCtx, filter, opts)
		if err != nil {
			return nil, err
		}

		err = cur.All(sessCtx, &questions)
		return questions, err
	})

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

func (service *QuestionBoxService) DeleteQuestion(questionID primitive.ObjectID) (sessErr error) {
	transaction := func(sessCtx mongo.SessionContext) (res interface{}, err error) {
		question := new(models.QuestionBoxQuestion)
		err = db.GetCollection("questions").FindOneAndDelete(sessCtx, bson.M{"_id": questionID}).Decode(question)
		if err != nil {
			return
		}

		filter := bson.M{
			"questions": bson.M{
				"$elemMatch": bson.M{
					"questionID": questionID,
				},
			},
		}

		update := bson.M{
			"$pull": bson.M{
				"questions": bson.M{"questionID": questionID},
			},
			"$inc": bson.M{
				"stats.questionNum": -1,
			},
		}
		_, err = db.GetCollection("labels").UpdateMany(sessCtx, filter, update)
		return
	}

	sessOpts := options.Session().SetDefaultWriteConcern(service.wc)
	_, sessErr = db.GetMongoConn().UseSession(sessOpts, transaction)
	return sessErr
}

/*
--------------------------------------------------
标签（文件夹）模块区域
*/

// NewLabels 创建多个标签
func (service *QuestionBoxService) NewLabels(labels []*models.QuestionLabel) (labelIDs []primitive.ObjectID, err error) {
	// 如果问题不存在标签，则直接退出
	if labels == nil {
		return
	}

	sessOpts := options.Session().SetDefaultWriteConcern(service.wc).SetDefaultReadConcern(service.rc)

	labelUserIDs := make([]primitive.ObjectID, 0)
	labelContents := make([]string, 0)

	for _, label := range labels {
		if label.Content == "" {
			err = fmt.Errorf("部分标签不存在内容")
			continue
		}

		if label.UserID == primitive.NilObjectID {
			err = fmt.Errorf("标签%s不存在用户信息", label.Content)
			continue
		}

		label.Init()

		labelUserIDs = append(labelUserIDs, label.UserID)
		labelContents = append(labelContents, label.Content)
	}

	filter := bson.M{
		"userID": bson.M{
			"$in": labelUserIDs,
		},
		"content": bson.M{
			"$in": labelContents,
		},
	}

	_, err = db.GetMongoConn().UseSession(sessOpts, func(sessCtx mongo.SessionContext) (interface{}, error) {
		cur, sessErr := db.GetCollection("labels").Find(sessCtx, filter)
		if sessErr != nil {
			return nil, sessErr
		}

		existsLabels := make([]*models.QuestionLabel, 0)
		sessErr = cur.All(sessCtx, &existsLabels)
		if sessErr != nil {
			return nil, sessErr
		}

		diffFunc := func(arrA, arrB []*models.QuestionLabel) (diffArr []interface{}) {
			// arrA - arrB，差集
			labelArrMap := make(map[string]bool)
			for _, label := range arrB {
				labelArrMap[label.UserID.String()+label.Content] = true
			}

			for _, label := range arrA {
				if !labelArrMap[label.UserID.String()+label.Content] {
					diffArr = append(diffArr, label)
				}
			}
			return
		}

		diffLabels := diffFunc(labels, existsLabels)
		if len(diffLabels) == 0 {
			return nil, nil
		}

		result, sessErr := db.GetCollection("labels").InsertMany(sessCtx, diffLabels)
		if sessErr != nil {
			return nil, sessErr
		}

		for _, labelID := range result.InsertedIDs {
			labelIDs = append(labelIDs, labelID.(primitive.ObjectID))
		}

		return nil, nil
	})

	return
}

func (service *QuestionBoxService) QueryLabelByID(labelID primitive.ObjectID) (label *models.QuestionLabel, err error) {
	label = new(models.QuestionLabel)
	err = db.GetCollection("labels").FindOne(context.TODO(), bson.M{"_id": labelID}).Decode(label)
	return
}

func (service *QuestionBoxService) QueryLabelByContent(content string) (label *models.QuestionLabel, err error) {
	label = new(models.QuestionLabel)
	err = db.GetCollection("labels").FindOne(context.TODO(), bson.M{"content": content}).Decode(label)
	return
}

func (service *QuestionBoxService) QueryLabelsFromUser(user *models.User, page, pageNum int64) (labels []*models.QuestionLabel, err error) {
	filter := bson.M{
		"userID": user.ID,
	}

	if page < 0 || pageNum <= 0 {
		err = fmt.Errorf("page: %d, pageNum: %d", page, pageNum)
		return
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

func (service *QuestionBoxService) QueryLabelsFromQuestion(user *models.User, question *models.QuestionBoxQuestion, page, pageNum int64) (labels []*models.QuestionLabel, err error) {
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

func (service *QuestionBoxService) DeleteLabel(labelID primitive.ObjectID) (err error) {
	err = db.GetCollection("labels").FindOneAndDelete(context.TODO(), bson.M{"_id": labelID}).Err()
	return
}

func (service *QuestionBoxService) UpdateLabelContent(label *models.QuestionLabel) (err error) {
	if label.Content == "" {
		err = errors.New("标签内容为空")
		return
	}

	if label.ID == primitive.NilObjectID {
		err = errors.New("标签用户ID为空")
		return
	}

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
		ID  primitive.ObjectID `bson:"_id"`
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
		"_id": labelID,
	}

	update := bson.M{
		"$pull": bson.M{
			"questions": bson.M{"questionID": questionID},
		},
		"$inc": bson.M{
			"stats.questionNum": -1,
		},
		"$set": bson.M{
			"updateTime": uint64(time.Now().Unix()),
		},
	}

	err = db.GetCollection("labels").FindOneAndUpdate(context.TODO(), filter, update).Err()
	return
}

func (service *QuestionBoxService) ChangeQuestionReadStatusInLabel(labelID primitive.ObjectID, questionID primitive.ObjectID) (err error) {

	filter := bson.M{
		"_id":                  labelID,
		"questions.questionID": questionID,
	}

	update := bson.M{
		"$set": bson.M{
			"questions.$.hasRead": true,
		},
	}

	type AggResult struct {
		ID       primitive.ObjectID         `bson:"_id"`
		Question models.QuestionInLabelInfo `bson:"questions"`
	}

	aggResult := make([]AggResult, 0)

	transaction := func(sessCtx mongo.SessionContext) (res interface{}, sessErr error) {
		cur, sessErr := db.GetCollection("labels").Aggregate(context.TODO(), mongo.Pipeline{
			{{"$match", bson.M{"_id": labelID}}},
			{{"$project", bson.M{"questions": 1}}},
			{{"$unwind", "$questions"}},
			{{"$match", bson.M{"questions.questionID": questionID}}},
		})
		if sessErr != nil {
			return
		}

		sessErr = cur.All(context.TODO(), &aggResult)
		if sessErr != nil {
			return
		}

		if aggResult[0].Question.HasRead {
			sessErr = errors.New("questionID对应问题状态已为已读")
			return
		}

		db.GetCollection("labels").FindOneAndUpdate(context.TODO(), filter, update)
		return
	}

	sessOpts := options.Session().SetDefaultWriteConcern(service.wc).SetDefaultReadConcern(service.rc)
	_, err = db.GetMongoConn().UseSession(sessOpts, transaction)
	return
}

func (service *QuestionBoxService) CountReadQuestionInLabel(labelID primitive.ObjectID) (questionReadNum int, err error) {
	sessOpts := options.Session().SetDefaultWriteConcern(service.wc).SetDefaultReadConcern(service.rc)

	type AggResult struct {
		QuestionNum int `bson:"questions"`
	}
	aggResult := make([]AggResult, 0)
	_, err = db.GetMongoConn().UseSession(sessOpts, func(sessCtx mongo.SessionContext) (result interface{}, sessErr error) {
		cur, sessErr := db.GetCollection("labels").Aggregate(sessCtx, mongo.Pipeline{
			{{"$match", bson.M{"_id": labelID}}},
			{{"$project", bson.M{"questions": 1}}},
			{{"$unwind", "$questions"}},
			{{"$match", bson.M{"questions.hasRead": true}}},
			{{"$count", "questions"}},
		})
		if sessErr != nil {
			return
		}

		sessErr = cur.All(sessCtx, &aggResult)
		return
	})
	if err != nil {
		return
	}

	if len(aggResult) == 0 {
		questionReadNum = 0
	} else {
		questionReadNum = aggResult[0].QuestionNum
	}

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
	/*TODO：这里是将answer和question关联，调用了一下师兄的AddAnswerToQuestion函数。

	这里InsertOne和AddAnswerToQuestion感觉类似一个事务，一个执行失败另一个也要回滚；比如
	这里如果已经执行InsertOne成功，但AddAnswerToQuestion执行失败，那这个回答即使存入数据库也是
	没有用的数据了（没有和question绑定，后续用questionID就查不出来）。

	如果这种没用的数据没有什么影响 或者这两个函数执行的成功率都很高以至于没用的数据很少，
	或许就这样写着？只让前端提示用户“提交回答失败”就行？

	resp: 不建议直接调用，我认为可以将我的那个函数删去，直接融入你的函数里面
	这里的原因是需要调用db.GetMongoConn().UseSession()函数，这个函数需要传递sessCtx这个上下文，
	就不能单纯的用context.TODO()或者context.BackGround()来解决了
	*/
	res, err := db.GetCollection("questionboxanswer").InsertOne(context.Background(), answer)
	if err != nil {
		return
	}
	docID = res.InsertedID.(primitive.ObjectID)
	err = questionBoxService.AddAnswerToQuestion(questionID, answer.ID)
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

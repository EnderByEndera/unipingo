package questionbox_test

import (
	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"testing"
	"time"
)

import (
	"context"
	"melodie-site/server/db"
	"strconv"
)

func getLabels(labelNum int) (labels []*models.QuestionBoxLabel, err error) {
	labels = make([]*models.QuestionBoxLabel, labelNum)
	user, err := services.GetAuthService().GetUserByName("admin")
	if err != nil {
		return
	}

	questions, err := services.GetQuestionBoxService().QueryQuestionsFromUser(&user, 0, 10)
	if err != nil {
		return
	}

	questionInLabelInfos := make([]models.QuestionInLabelInfo, len(questions))
	for index := range questionInLabelInfos {
		questionInLabelInfos[index] = models.QuestionInLabelInfo{
			ID:   questions[index].ID,
			Name: questions[index].Title,
		}
	}

	for index := range labels {
		labels[index] = new(models.QuestionBoxLabel)
		labels[index].ID = primitive.NewObjectID()
		labels[index].Content = "Hello World " + primitive.NewObjectID().String()
		labels[index].UserID = user.ID
		labels[index].Questions = questionInLabelInfos
	}

	return
}

func TestNewLabels(t *testing.T) {
	labels, err := getLabels(10)
	assert.Equal(t, err, nil)

	documentNumBefore, err := db.GetCollection("labels").CountDocuments(context.TODO(), bson.M{})
	assert.Equal(t, err, nil)

	labelIDs, err := services.GetQuestionBoxService().NewLabels(labels)
	assert.Equal(t, err, nil)

	documentNum, err := db.GetCollection("labels").CountDocuments(context.TODO(), bson.M{})
	assert.Equal(t, err, nil)
	assert.Equal(t, documentNum, documentNumBefore+10)

	labelIDs, err = services.GetQuestionBoxService().NewLabels(labels)
	assert.Equal(t, err, nil)

	defer func() {
		for _, labelID := range labelIDs {
			_ = services.GetQuestionBoxService().DeleteLabel(labelID)
		}
	}()

	documentNumReal, err := db.GetCollection("labels").CountDocuments(context.TODO(), bson.M{})
	assert.Equal(t, err, nil)

	assert.Equal(t, documentNum, documentNumReal)

	label := labels[0]
	label.Questions = []models.QuestionInLabelInfo{
		{
			ID:   primitive.NewObjectID(),
			Name: "Test " + primitive.NewObjectID().String(),
		},
		{
			ID:   primitive.NewObjectID(),
			Name: "Test " + primitive.NewObjectID().String(),
		},
	}

	labels, err = getLabels(9)
	assert.Equal(t, err, nil)
	// 此时labels中有9个新label，一个替换了问题的原label
	labels = append(labels, label)

	labelIDs, err = services.GetQuestionBoxService().NewLabels(labels)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(labelIDs), 9)

	updatedLabel := new(models.QuestionBoxLabel)
	err = db.GetCollection("labels").FindOne(context.TODO(), bson.M{"_id": label.ID}).Decode(updatedLabel)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(updatedLabel.Questions), 12) // 12是因为getLabels()中初始化问题数量为10，再加上新增的两个test问题一共12
}

func TestQueryLabelFromUser(t *testing.T) {
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)

	labels, err := services.GetQuestionBoxService().QueryLabelsFromUser(&user, 0, 10)
	assert.Equal(t, err, nil)
	for _, label := range labels {
		assert.Equal(t, label.UserID, user.ID)
	}

	_, err = services.GetQuestionBoxService().QueryLabelsFromUser(&user, 0, 0)
	assert.NotEqual(t, err, nil)

	_, err = services.GetQuestionBoxService().QueryLabelsFromUser(&user, -1, 0)
	assert.NotEqual(t, err, nil)

	var page int64 = 2
	var pageNum int64 = 626737562 // 尝试较大的数
	_, err = services.GetQuestionBoxService().QueryLabelsFromUser(&user, page, pageNum)
	assert.Equal(t, err, nil)
}

func TestQueryLabelFromQuestion(t *testing.T) {
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)

	questions, err := services.GetQuestionBoxService().QueryQuestionsFromUser(&user, 0, 1)
	assert.Equal(t, err, nil)

	labels, err := services.GetQuestionBoxService().QueryLabelsFromUser(&user, 0, 100)
	assert.Equal(t, err, nil)

	for _, label := range labels {
		for _, question := range questions {
			questionInfo := &models.QuestionInLabelInfo{ID: question.ID, Name: question.Title}
			err := services.GetQuestionBoxService().AddQuestionInLabel(label.ID, questionInfo)
			assert.Equal(t, err, nil)
		}
	}

	for _, question := range questions {
		sLabels, err := services.GetQuestionBoxService().QueryLabelsFromQuestion(&user, question, 0, 100)
		assert.Equal(t, err, nil)
		assert.NotEqual(t, len(sLabels), 0)
		assert.Equal(t, len(labels), len(sLabels))
	}
}

func TestUpdateLabelContent(t *testing.T) {
	labels, err := getLabels(10)
	assert.Equal(t, err, nil)
	labelIDs, err := services.GetQuestionBoxService().NewLabels(labels)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, len(labelIDs), 0)

	defer func() {
		for _, id := range labelIDs {
			_ = services.GetQuestionBoxService().DeleteLabel(id)
		}
	}()

	newLabels := make([]*models.QuestionBoxLabel, 0)
	for _, labelID := range labelIDs {
		// 如果要更新内容
		label := &models.QuestionBoxLabel{
			ID:      labelID,
			Content: "Hello From World",
		}
		newLabels = append(newLabels, label)
	}

	for _, label := range newLabels {
		err = services.GetQuestionBoxService().UpdateLabelContent(label)
		assert.Equal(t, err, nil)

		sLabel, err := services.GetQuestionBoxService().QueryLabelByContent(label.Content)
		assert.Equal(t, err, nil)
		assert.NotEqual(t, sLabel, nil)
	}
}

func TestAddQuestionInLabel(t *testing.T) {
	question, err := getOneQuestion("TestAddQuestionInLabel")
	assert.Equal(t, err, nil)
	questionID, err := services.GetQuestionBoxService().NewQuestion(question)
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
	}()
	assert.Equal(t, err, nil)

	questionInfo := &models.QuestionInLabelInfo{
		ID:   questionID,
		Name: question.Title,
	}

	labels, err := getLabels(1)
	assert.Equal(t, err, nil)
	labelIDs, err := services.GetQuestionBoxService().NewLabels(labels)

	defer func() {
		for _, id := range labelIDs {
			_ = services.GetQuestionBoxService().DeleteLabel(id)
		}
	}()
	assert.Equal(t, err, nil)

	newLabels := make([]*models.QuestionBoxLabel, 0)
	for index, id := range labelIDs {
		label, err := services.GetQuestionBoxService().QueryLabelByID(id)
		assert.Equal(t, err, nil)
		newLabels = append(newLabels, label)

		err = services.GetQuestionBoxService().AddQuestionInLabel(id, questionInfo)
		assert.Equal(t, err, nil)

		addedLabel, err := services.GetQuestionBoxService().QueryLabelByID(labelIDs[index])
		assert.Equal(t, err, nil)
		assert.Equal(t, addedLabel.Content, newLabels[index].Content)
	}

}

func TestDeleteQuestionFromLabel(t *testing.T) {
	deletedQuestion, err := getOneQuestion("TestDeleteQuestionFromLabel")
	assert.Equal(t, err, nil)
	questionID, _ := services.GetQuestionBoxService().NewQuestion(deletedQuestion)
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
	}()

	deletedQuestionInfo := &models.QuestionInLabelInfo{
		ID:   deletedQuestion.ID,
		Name: deletedQuestion.Title,
	}

	labels, err := getLabels(1)
	labelIDs, err := services.GetQuestionBoxService().NewLabels(labels)
	assert.Equal(t, err, nil)
	defer func() {
		for _, id := range labelIDs {
			_ = services.GetQuestionBoxService().DeleteLabel(id)
		}
	}()

	assert.NotEqual(t, len(labelIDs), 0)

	for index := range labelIDs {
		err = services.GetQuestionBoxService().AddQuestionInLabel(labelIDs[index], deletedQuestionInfo)
		assert.Equal(t, err, nil)

		addedLabel, err := services.GetQuestionBoxService().QueryLabelByID(labelIDs[index])
		assert.Equal(t, err, nil)

		err = services.GetQuestionBoxService().DeleteQuestionFromLabel(addedLabel.ID, deletedQuestionInfo.ID)
		assert.Equal(t, err, nil)

		deletedLabel, err := services.GetQuestionBoxService().QueryLabelByID(labelIDs[index])
		assert.Equal(t, err, nil)

		for _, question := range deletedLabel.Questions {
			assert.NotEqual(t, question.ID, deletedQuestion.ID)
		}
	}
}

func TestQuestionHasReadInLabel(t *testing.T) {
	//TODO: Rewrite
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)

	labels, err := getLabels(1)
	assert.Equal(t, err, nil)

	question, err := getOneQuestion("TestQuestionHasReadInLabel")
	assert.Equal(t, err, nil)

	questionID, err := services.GetQuestionBoxService().NewQuestion(question)
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
	}()
	assert.Equal(t, err, nil)

	for index := range labels {
		labels[index].Questions = append(labels[index].Questions, models.QuestionInLabelInfo{
			ID:   question.ID,
			Name: question.Title,
		})
	}
	assert.Equal(t, err, nil)
	labelIDs, err := services.GetQuestionBoxService().NewLabels(labels)
	defer func() {
		for _, labelID := range labelIDs {
			_ = services.GetQuestionBoxService().DeleteLabel(labelID)
		}
	}()
	assert.Equal(t, err, nil)

	answers := make([]*models.QuestionBoxAnswer, 0)
	answerIDs := make([]primitive.ObjectID, 0)

	for i := 0; i < 10; i++ {
		answer, err := getOneAnswer(question)
		assert.Equal(t, err, nil)
		answer.QuestionID = question.ID
		answerID, err := services.GetQuestionBoxService().NewAnswer(answer)
		answerIDs = append(answerIDs, answerID)
		assert.Equal(t, err, nil)
		question.Answers = append(question.Answers, answerID)
		answers = append(answers, answer)
	}

	defer func() {
		for _, answerID := range answerIDs {
			_ = services.GetQuestionBoxService().DeleteQuestionBoxAnswerByID(answerID)
		}
	}()

	for index := range answers {
		err = services.GetQuestionBoxService().ReadAnswerByUser(user.ID, answers[index])
		assert.Equal(t, err, nil)
		readNum, err := services.GetQuestionBoxService().CountAnswerReadNumInQuestion(question.ID, user.ID)
		assert.Equal(t, err, nil)
		assert.Equal(t, readNum, index+1)
	}
}

func BenchmarkNewLabels(b *testing.B) {
	labels, err := getLabels(100)
	assert.Equal(b, err, nil)

	b.ResetTimer()
	b.SetParallelism(36)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			labelIDs, err := services.GetQuestionBoxService().NewLabels(labels)
			assert.Equal(b, err, nil)
			println(labelIDs)
		}
	})
}

func BenchmarkQueryLabelFromUser(b *testing.B) {
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(b, err, nil)

	b.ResetTimer()
	b.SetParallelism(3)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := services.GetQuestionBoxService().QueryLabelsFromUser(&user, int64(rand.Int()%1000+1), int64(rand.Int()%1000+1))
			assert.Equal(b, err, nil)
		}
	})
}

func BenchmarkQueryLabelFromQuestion(b *testing.B) {
	question, err := getOneQuestion("BenchmarkQueryLabelFromQuestion")
	assert.Equal(b, err, nil)
	questionID, err := services.GetQuestionBoxService().NewQuestion(question)
	assert.Equal(b, err, nil)
	question.ID = questionID
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
	}()

	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(b, err, nil)

	labelIDs, err := services.GetQuestionBoxService().QueryLabelsFromUser(&user, 0, 1000)

	for _, labelID := range labelIDs {
		questionInfo := &models.QuestionInLabelInfo{
			ID:   question.ID,
			Name: question.Title,
		}
		err := services.GetQuestionBoxService().AddQuestionInLabel(labelID.ID, questionInfo)
		assert.Equal(b, err, nil)
	}

	b.ResetTimer()
	b.SetParallelism(3)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().Unix())
			_, err = services.GetQuestionBoxService().QueryLabelsFromQuestion(&user, question, int64(rand.Int()%1000+1), int64(rand.Int()%1000+1))
			assert.Equal(b, err, nil)
		}
	})
}

func BenchmarkDeleteQuestionFromLabel(b *testing.B) {
	questionIDs := make([]primitive.ObjectID, 0)

	for i := 0; i < 10; i++ {
		rand.Seed(time.Now().Unix() + int64(i))
		question, err := getOneQuestion("BenchmarkDeleteQuestionFromLabel " + strconv.Itoa(rand.Int()))
		assert.Equal(b, err, nil)

		questionID, err := services.GetQuestionBoxService().NewQuestion(question)
		assert.Equal(b, err, nil)

		questionIDs = append(questionIDs, questionID)
	}

	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(b, err, nil)

	labelIDs, err := services.GetQuestionBoxService().QueryLabelsFromUser(&user, 0, 10000)
	assert.Equal(b, err, nil)

	for _, labelID := range labelIDs {
		for _, questionID := range questionIDs {
			questionInfo := &models.QuestionInLabelInfo{
				ID:   questionID,
				Name: "BenchmarkDeleteQuestionFromLabel",
			}
			err := services.GetQuestionBoxService().AddQuestionInLabel(labelID.ID, questionInfo)
			assert.Equal(b, err, nil)
		}
	}

	b.ResetTimer()
	b.StartTimer()
	for _, questionID := range questionIDs {
		_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
	}
	b.StopTimer()
}

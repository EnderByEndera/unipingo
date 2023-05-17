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

func getLabels(labelNum int) (labels []*models.QuestionLabel, err error) {
	labels = make([]*models.QuestionLabel, labelNum)
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
			ID:      questions[index].ID,
			Name:    questions[index].Title,
			HasRead: false,
		}
	}

	for index := range labels {
		labels[index] = new(models.QuestionLabel)
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

	documentNumTwo, err := db.GetCollection("labels").CountDocuments(context.TODO(), bson.M{})
	assert.Equal(t, err, nil)

	assert.Equal(t, documentNum, documentNumTwo)
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

	var page int64 = 0
	var pageNum int64 = 626737562
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
			questionInfo := &models.QuestionInLabelInfo{ID: question.ID, Name: question.Title, HasRead: false}
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

	newLabels := make([]*models.QuestionLabel, 0)
	for _, labelID := range labelIDs {
		// 如果要更新内容
		label := &models.QuestionLabel{
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
		ID:      questionID,
		Name:    question.Title,
		HasRead: false,
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

	newLabels := make([]*models.QuestionLabel, 0)
	for index, id := range labelIDs {
		label, err := services.GetQuestionBoxService().QueryLabelByID(id)
		assert.Equal(t, err, nil)
		newLabels = append(newLabels, label)

		err = services.GetQuestionBoxService().AddQuestionInLabel(id, questionInfo)
		assert.Equal(t, err, nil)

		addedLabel, err := services.GetQuestionBoxService().QueryLabelByID(labelIDs[index])
		assert.Equal(t, err, nil)
		assert.Equal(t, addedLabel.Content, newLabels[index].Content)
		assert.Equal(t, addedLabel.Statistics.QuestionNum, newLabels[index].Statistics.QuestionNum+1)
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
		ID:      deletedQuestion.ID,
		Name:    deletedQuestion.Title,
		HasRead: false,
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
		assert.Equal(t, deletedLabel.Statistics.QuestionNum, addedLabel.Statistics.QuestionNum-1)
	}
}

func TestQuestionHasReadInLabel(t *testing.T) {
	labels, err := getLabels(1)
	assert.Equal(t, err, nil)

	questionInfos := make([]models.QuestionInLabelInfo, 0)
	for i := 0; i < 10; i++ {
		question, err := getOneQuestion(primitive.NewObjectID().String())
		assert.Equal(t, err, nil)
		newQuestionID, err := services.GetQuestionBoxService().NewQuestion(question)
		defer func() {
			_ = services.GetQuestionBoxService().DeleteQuestion(newQuestionID)
		}()
		assert.Equal(t, err, nil)

		questionInfos = append(questionInfos, models.QuestionInLabelInfo{
			ID:      newQuestionID,
			Name:    question.Title,
			HasRead: false,
		})
	}

	for index := range labels {
		labels[index].Questions = questionInfos
	}

	labelIDs, err := services.GetQuestionBoxService().NewLabels(labels)
	defer func() {
		for index := range labelIDs {
			_ = services.GetQuestionBoxService().DeleteLabel(labelIDs[index])
		}
	}()
	assert.Equal(t, err, nil)

	newLabels := make([]*models.QuestionLabel, 0)
	for index := range labelIDs {
		label, err := services.GetQuestionBoxService().QueryLabelByID(labelIDs[index])
		assert.Equal(t, err, nil)

		newLabels = append(newLabels, label)
	}

	for index := range labelIDs {
		err = services.GetQuestionBoxService().ChangeQuestionReadStatusInLabel(labelIDs[index], newLabels[index].Questions[0].ID)
		assert.Equal(t, err, nil)

		label, err := services.GetQuestionBoxService().QueryLabelByID(labelIDs[index])
		assert.Equal(t, err, nil)

		ok := false
		for _, question := range label.Questions {
			if question.HasRead {
				ok = true
				break
			}
		}
		assert.Equal(t, ok, true)
	}

}

func TestCountReadQuestionInLabel(t *testing.T) {
	labels, err := getLabels(1)
	assert.Equal(t, err, nil)

	labelIDs, err := services.GetQuestionBoxService().NewLabels(labels)
	assert.Equal(t, err, nil)

	defer func() {
		for _, labelID := range labelIDs {
			_ = services.GetQuestionBoxService().DeleteLabel(labelID)
		}
	}()

	newLabels := make([]*models.QuestionLabel, 0)
	for _, labelID := range labelIDs {
		label, err := services.GetQuestionBoxService().QueryLabelByID(labelID)
		assert.Equal(t, err, nil)
		newLabels = append(newLabels, label)
	}

	for _, newLabel := range newLabels {
		count := 0

		for index := range newLabel.Questions {
			if rand.Intn(2) == 0 {
				newLabel.Questions[index].HasRead = true
				count++
			}
		}
		err = db.GetCollection("labels").FindOneAndUpdate(context.TODO(),
			bson.M{"_id": newLabel.ID},
			bson.M{"$set": bson.M{"questions": newLabel.Questions}}).Err()
		assert.Equal(t, err, nil)

		readNum, err := services.GetQuestionBoxService().CountReadQuestionInLabel(newLabel.ID)
		assert.Equal(t, err, nil)
		assert.Equal(t, readNum, count)
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
			ID:      question.ID,
			Name:    question.Title,
			HasRead: false,
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
				ID:      questionID,
				Name:    "BenchmarkDeleteQuestionFromLabel",
				HasRead: false,
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

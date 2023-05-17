package questionbox_test

import (
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/jinzhu/copier"
	"math/rand"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"strconv"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getOneQuestion(title string) (*models.QuestionBoxQuestion, error) {
	admin, err := services.GetAuthService().GetUserByName("admin")
	if err != nil {
		return nil, err
	}

	school, err := services.GetHEIService().GetHEIByName("中国科学院大学")
	if err != nil {
		return nil, err
	}

	major, err := services.GetMajorService().GetMajorByName("哲学")
	if err != nil {
		return nil, err
	}

	question := &models.QuestionBoxQuestion{
		ID:          primitive.NewObjectID(),
		UserID:      admin.ID,
		Title:       title,
		Description: "This is a new question",
		School: models.EntityWithName{
			ID:   school.ID,
			Name: school.Name,
		},
		Major: models.EntityWithName{
			ID:   major.ID,
			Name: major.Name,
		},
		Questioner: models.PersonalInfo{
			CEEPlace:  "北京",
			Subject:   "语文",
			Age:       17,
			Gender:    "男",
			Situation: "",
		},
	}
	return question, nil
}

func TestNewQuestion(t *testing.T) {
	question, err := getOneQuestion("TestNewQuestion")
	assert.Equal(t, err, nil)

	docID, err := services.GetQuestionBoxService().NewQuestion(question)
	deferDocID := docID
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(deferDocID)
	}()
	assert.Equal(t, err, nil)
	assert.NotEqual(t, docID, primitive.NilObjectID)

	docID, err = services.GetQuestionBoxService().NewQuestion(question)
	assert.Equal(t, err, errors.New("该问题已存在"))
	assert.NotEqual(t, docID, primitive.NilObjectID)

	questionWithoutDescription := new(models.QuestionBoxQuestion)
	err = copier.Copy(questionWithoutDescription, question)
	assert.Equal(t, err, nil)

	questionWithoutDescription.Title = ""
	questionWithoutDescription.Description = ""
	docID, err = services.GetQuestionBoxService().NewQuestion(questionWithoutDescription)
	assert.Equal(t, err, errors.New("该问题没有填写标题或描述"))
	assert.Equal(t, docID, primitive.NilObjectID)

	questionWithoutMajorAndSchool := new(models.QuestionBoxQuestion)
	err = copier.Copy(questionWithoutMajorAndSchool, question)
	assert.Equal(t, err, nil)

	questionWithoutMajorAndSchool.Major = models.EntityWithName{}
	questionWithoutMajorAndSchool.School = models.EntityWithName{}
	docID, err = services.GetQuestionBoxService().NewQuestion(questionWithoutMajorAndSchool)
	assert.Equal(t, err, errors.New("该问题学校和专业均为空"))
	assert.Equal(t, docID, primitive.NilObjectID)
}

func TestNewQuestionWithExists(t *testing.T) {
	question, err := getOneQuestion("My Question")
	assert.Equal(t, err, nil)
	docID, err := services.GetQuestionBoxService().NewQuestion(question)
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(docID)
	}()
	assert.Equal(t, err, nil)

	docIDSec, err := services.GetQuestionBoxService().NewQuestion(question)
	assert.NotEqual(t, err, nil)
	assert.Equal(t, docID, docIDSec)
}

func BenchmarkNewQuestion(b *testing.B) {
	questions := make([]*models.QuestionBoxQuestion, 10000)
	for i := range questions {
		question, _ := getOneQuestion(strconv.Itoa(i))
		questions[i] = question
	}


	b.ResetTimer()
	b.SetParallelism(36)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			rand.Seed(time.Now().Unix())
			_, _ = services.GetQuestionBoxService().NewQuestion(questions[rand.Int()%10000])
		}
	})

	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(b, err, nil)
	questions, err = services.GetQuestionBoxService().QueryQuestionsFromUser(&user, 0, 20)
	assert.Equal(b, err, nil)
}

func TestGetQuestion(t *testing.T) {
	question, err := getOneQuestion("GetQuestion")
	assert.Equal(t, err, nil)
	questionID, err := services.GetQuestionBoxService().NewQuestion(question)
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
	}()
	assert.Equal(t, err, nil)

	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)
	questions, err := services.GetQuestionBoxService().QueryQuestionsFromUser(&user, 0, 1000)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, len(questions), 0)
	for index := range questions {
		question, err := services.GetQuestionBoxService().QueryQuestionByID(questions[index].ID)
		assert.Equal(t, err, nil)
		assert.NotEqual(t, question.CreateTime, nil)
	}
}

func BenchmarkGetQuestion(b *testing.B) {
	questionID, err := primitive.ObjectIDFromHex("64536b800a3da36ef0a12770")
	assert.Equal(b, err, nil)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			question, err := services.GetQuestionBoxService().QueryQuestionByID(questionID)
			assert.Equal(b, err, nil)
			assert.NotEqual(b, question.ID, primitive.NilObjectID)
		}
	})
}

func TestQuestionList(t *testing.T) {
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)
	_, err = services.GetQuestionBoxService().QueryQuestionsFromUser(&user, 0, 10)
	assert.Equal(t, err, nil)
	_, err = services.GetQuestionBoxService().QueryQuestionsFromUser(&user, -1, 10)
	assert.NotEqual(t, err, nil)
	_, err = services.GetQuestionBoxService().QueryQuestionsFromUser(&user, -1, -10)
	assert.NotEqual(t, err, nil)
	_, err = services.GetQuestionBoxService().QueryQuestionsFromUser(&user, 1, -10)
	assert.NotEqual(t, err, nil)
	_, err = services.GetQuestionBoxService().QueryQuestionsFromUser(nil, 1, 10)
	assert.NotEqual(t, err, nil)
}

func BenchmarkQuestionList(b *testing.B) {
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(b, err, nil)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, _ = services.GetQuestionBoxService().QueryQuestionsFromUser(&user, int64(rand.Int()), int64(rand.Int()))
		}
	})
}

func TestUpdateQuestion(t *testing.T) {
	question, err := getOneQuestion("My New Question")
	assert.Equal(t, err, nil)
	questionID, err := services.GetQuestionBoxService().NewQuestion(question)

	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
	}()

	assert.Equal(t, err, nil)
	newQuestion := &models.QuestionBoxQuestion{
		ID:          questionID,
		Description: "This is a new description",
	}
	err = services.GetQuestionBoxService().UpdateQuestionDescription(newQuestion)
	assert.Equal(t, err, nil)
	newQuestion.Description = ""
	err = services.GetQuestionBoxService().UpdateQuestionDescription(newQuestion)
	assert.Equal(t, err, errors.New("更新描述为空"))
	newQuestion = &models.QuestionBoxQuestion{
		ID:          primitive.NilObjectID,
		Description: "This is a new description",
	}
	err = services.GetQuestionBoxService().UpdateQuestionDescription(newQuestion)
	assert.NotEqual(t, err, nil)

	_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
}

func BenchmarkUpdateQuestion(b *testing.B) {
	question, err := getOneQuestion("My New Question")
	questionID, err := services.GetQuestionBoxService().NewQuestion(question)
	assert.Equal(b, err, nil)
	newQuestion := &models.QuestionBoxQuestion{
		ID:          questionID,
		Description: "This is a new description",
	}

	b.SetParallelism(100)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			err = services.GetQuestionBoxService().UpdateQuestionDescription(newQuestion)
			assert.Equal(b, err, nil)
		}
	})

	_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
}

func TestDeleteQuestion(t *testing.T) {
	question, err := getOneQuestion("Delete Question")
	assert.Equal(t, err, nil)

	questionInfo := &models.QuestionInLabelInfo{
		ID:      question.ID,
		Name:    question.Title,
		HasRead: false,
	}

	questionID, err := services.GetQuestionBoxService().NewQuestion(question)

	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
	}()
	assert.Equal(t, err, nil)

	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
	}()

	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)

	labels, err := services.GetQuestionBoxService().QueryLabelsFromUser(&user, 0, 10)
	assert.Equal(t, err, nil)

	labelIDs := make([]primitive.ObjectID, 0)
	for _, label := range labels {
		labelIDs = append(labelIDs, label.ID)
	}

	for _, labelID := range labelIDs {
		err = services.GetQuestionBoxService().AddQuestionInLabel(labelID, questionInfo)
		assert.Equal(t, err, nil)
	}

	err = services.GetQuestionBoxService().DeleteQuestion(questionID)
	assert.Equal(t, err, nil)

	err = services.GetQuestionBoxService().DeleteQuestion(questionID)
	assert.NotEqual(t, err, nil)
}

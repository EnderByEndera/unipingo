package tests

import (
	"errors"
	"math/rand"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"strconv"
	"testing"

	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getOneQuestion(t *testing.T) *models.QuestionBoxQuestion {
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)
	school, err := services.GetHEIService().GetHEIByName("中国科学院大学")
	assert.Equal(t, err, nil)
	major, err := services.GetMajorService().GetMajorByName("哲学")
	assert.Equal(t, err, nil)
	question := &models.QuestionBoxQuestion{
		UserID:      user.ID,
		Title:       "My Question",
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
	return question
}

func TestNewQuestion(t *testing.T) {
	question := getOneQuestion(t)
	docID, err := services.GetQuestionBoxService().NewQuestion(question)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, docID, primitive.NilObjectID)
	question.Title = ""
	question.Description = ""
	docID, err = services.GetQuestionBoxService().NewQuestion(question)
	assert.Equal(t, err, errors.New("该问题没有填写标题或描述"))
	assert.Equal(t, docID, primitive.NilObjectID)
}

func TestNewQuestionWithExists(t *testing.T) {
	question := getOneQuestion(t)
	_, _ = services.GetQuestionBoxService().NewQuestion(question)
	docID, err := services.GetQuestionBoxService().NewQuestion(question)
	assert.NotEqual(t, err, nil)
	assert.Equal(t, docID, primitive.NilObjectID)
}

func BenchmarkNewQuestion(b *testing.B) {
	questions := make([]*models.QuestionBoxQuestion, 10000)
	for i := range questions {
		question := &models.QuestionBoxQuestion{
			Description: "111",
		}
		question.Title = strconv.Itoa(i)
		questions[i] = question
	}
	// TODO  这里为什么要设定一下36呢，然后后边bench测试时有的没有设置这个值
	b.SetParallelism(36)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			docID, err := services.GetQuestionBoxService().NewQuestion(questions[rand.Int()%10000])
			if docID == primitive.NilObjectID {
				assert.Equal(b, err, errors.New("该问题已存在"))
			} else {
				assert.Equal(b, err, nil)
			}
		}
	})
	/*TODO
	  1、我这里想的是前面创建question时属性应该要传递userID，这样后面QuestionList才能查出来（不然这里QuestionList返回的应该是空数组吧?）
	  2、我运行一次，数据库question增加1000，为什么不是设定的是10000呢？
	*/
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(b, err, nil)
	questions, err = services.GetQuestionBoxService().QuestionList(&user, 0, 20)
	assert.Equal(b, err, nil)
}

func TestGetQuestion(t *testing.T) {
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)
	questions, err := services.GetQuestionBoxService().QuestionList(&user, 0, 1000)
	assert.Equal(t, err, nil)
	question, err := services.GetQuestionBoxService().QueryQuestionByID(questions[0].ID)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, question.CreateTime, nil)
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
	_, err = services.GetQuestionBoxService().QuestionList(&user, 0, 10)
	assert.Equal(t, err, nil)
	_, err = services.GetQuestionBoxService().QuestionList(&user, -1, 10)
	assert.NotEqual(t, err, nil)
	_, err = services.GetQuestionBoxService().QuestionList(&user, -1, -10)
	assert.NotEqual(t, err, nil)
	_, err = services.GetQuestionBoxService().QuestionList(&user, 1, -10)
	assert.NotEqual(t, err, nil)
	_, err = services.GetQuestionBoxService().QuestionList(nil, 1, 10)
	assert.NotEqual(t, err, nil)
}

func BenchmarkQuestionList(b *testing.B) {
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(b, err, nil)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			services.GetQuestionBoxService().QuestionList(&user, int64(rand.Int()), int64(rand.Int()))
		}
	})
}

func TestUpdateQuestion(t *testing.T) {
	questionID, err := primitive.ObjectIDFromHex("6453647c6fe2811ed76a9c04")
	assert.Equal(t, err, nil)
	question := &models.QuestionBoxQuestion{
		ID:          questionID,
		Description: "This is a new description",
	}
	err = services.GetQuestionBoxService().UpdateQuestionDescription(question)
	assert.Equal(t, err, nil)
	question.Description = ""
	err = services.GetQuestionBoxService().UpdateQuestionDescription(question)
	assert.Equal(t, err, errors.New("更新描述为空"))
	question = &models.QuestionBoxQuestion{
		ID:          primitive.NilObjectID,
		Description: "This is a new description",
	}
	err = services.GetQuestionBoxService().UpdateQuestionDescription(question)
	assert.NotEqual(t, err, nil)
}

func BenchmarkUpdateQuestion(b *testing.B) {
	questionID, err := primitive.ObjectIDFromHex("6453647c6fe2811ed76a9c04")
	assert.Equal(b, err, nil)
	question := &models.QuestionBoxQuestion{
		ID:          questionID,
		Description: "This is a new description",
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			err = services.GetQuestionBoxService().UpdateQuestionDescription(question)
			assert.Equal(b, err, nil)
		}
	})
}

func TestDeleteQuestion(t *testing.T) {
	questionID, err := primitive.ObjectIDFromHex("6460db3791ffc5d0aa7c26c2")
	assert.Equal(t, err, nil)

	err = services.GetQuestionBoxService().DeleteQuestion(questionID)
	assert.Equal(t, err, nil)

	err = services.GetQuestionBoxService().DeleteQuestion(questionID)
	assert.NotEqual(t, err, nil)
}

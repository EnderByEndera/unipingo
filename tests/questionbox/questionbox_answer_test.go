package questionbox_test

import (
	"errors"
	"math/rand"
	"melodie-site/server/models"
	"melodie-site/server/services"

	"testing"
	"time"

	"github.com/go-playground/assert/v2"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getOneAnswer(question *models.QuestionBoxQuestion) (answer *models.QuestionBoxAnswer, err error) {
	admin, err := services.GetAuthService().GetUserByName("admin")
	if err != nil {
		return
	}
	school, err := services.GetHEIService().GetHEIByName("华东师范大学")
	if err != nil {
		return
	}
	major, err := services.GetMajorService().GetMajorByName("软件工程")
	if err != nil {
		return
	}

	answer = &models.QuestionBoxAnswer{
		ID:         primitive.NewObjectID(),
		UserID:     admin.ID,
		CreateTime: uint64(time.Now().UnixMilli()),
		Content:    "华东师范大学中山北路校区在上海市普陀区中山北路3663号，学校很小，方便赶早八..",
		QuestionID: question.ID,
		School: models.EntityWithName{
			ID:   school.ID,
			Name: school.Name,
		},
		Major: models.EntityWithName{
			ID:   major.ID,
			Name: major.Name,
		},
		Respondant: models.PersonalInfo{
			CEEPlace:  "河南省",
			Subject:   "理科",
			Age:       20,
			Gender:    "女",
			Situation: "",
		},
	}

	return
}

func TestNewQuestionboxAnswer(t *testing.T) {
	question, err := getOneQuestion("TestNewQuestionboxAnswer")
	questionID, err := services.GetQuestionBoxService().NewQuestion(question)
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
	}()
	assert.Equal(t, err, nil)

	answer, err := getOneAnswer(question)
	assert.Equal(t, err, nil)

	answerID, err := services.GetQuestionBoxService().NewAnswer(answer)
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestionBoxAnswerByID(answerID)
	}()
	assert.Equal(t, err, nil)
	assert.NotEqual(t, answerID, primitive.NilObjectID)

	answer.Content = ""
	answerID, err = services.GetQuestionBoxService().NewAnswer(answer)
	assert.Equal(t, err, errors.New("该回答没有填写内容"))
	assert.Equal(t, answerID, primitive.NilObjectID)

	answer.Content = "认准华师大"
	answer.QuestionID, _ = primitive.ObjectIDFromHex("1")
	answerID, err = services.GetQuestionBoxService().NewAnswer(answer)
	assert.NotEqual(t, err, nil)
	assert.Equal(t, answerID, primitive.NilObjectID)
}

func BenchmarkNewQuestionboxAnswer(b *testing.B) {
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(b, err, nil)

	question, err := getOneQuestion(primitive.NewObjectID().String())
	assert.Equal(b, err, nil)
	questionID, err := services.GetQuestionBoxService().NewQuestion(question)
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
	}()
	assert.Equal(b, err, nil)

	answers := make([]*models.QuestionBoxAnswer, 1000)
	for i := range answers {
		answer := &models.QuestionBoxAnswer{
			QuestionID: questionID,
			UserID:     user.ID,
			Content:    "这里是回答的内容",
		}
		answers[i] = answer
	}

	b.SetParallelism(36)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := services.GetQuestionBoxService().NewAnswer(answers[rand.Int()%1000])
			assert.Equal(b, err, nil)
		}
	})

	answers, err = services.GetQuestionBoxService().AnswerList(question, 0, 20)
	assert.Equal(b, err, nil)
	answers, err = services.GetQuestionBoxService().MyAnswerList(&user, 0, 20)
	assert.Equal(b, err, nil)
}

func TestQueryAnswerByID(t *testing.T) {
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)
	answers, err := services.GetQuestionBoxService().MyAnswerList(&user, 0, 20)
	assert.Equal(t, err, nil)
	answer, err := services.GetQuestionBoxService().QueryAnswerByID(answers[0].ID)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, answer.CreateTime, nil)
}

func BenchmarkGetAnswer(b *testing.B) {
	//这里记得要补全
	answerID, err := primitive.ObjectIDFromHex("645df9d6d4145cfbbc2db78e")
	assert.Equal(b, err, nil)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			answer, err := services.GetQuestionBoxService().QueryAnswerByID(answerID)
			assert.Equal(b, err, nil)
			assert.NotEqual(b, answer.ID, primitive.NilObjectID)
		}
	})
}
func TestDeleteQuestionBoxAnswerByID(t *testing.T) {
	question, err := getOneQuestion("TestDeleteQuestionBoxAnswerByID")
	assert.Equal(t, err, nil)

	questionID, err := services.GetQuestionBoxService().NewQuestion(question)
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
	}()
	assert.Equal(t, err, nil)

	answer, err := getOneAnswer(question)
	assert.Equal(t, err, nil)

	answerID, err := services.GetQuestionBoxService().NewAnswer(answer)
	assert.Equal(t, err, nil)

	err = services.GetQuestionBoxService().DeleteQuestionBoxAnswerByID(answerID)
	assert.Equal(t, err, nil)
	err = services.GetQuestionBoxService().DeleteQuestionBoxAnswerByID(answerID)
	assert.NotEqual(t, err, nil)
}

func TestGetAnswerList(t *testing.T) {
	Questionid, _ := primitive.ObjectIDFromHex("645df98ee933a87982169266")
	question := &models.QuestionBoxQuestion{
		ID: Questionid,
	}
	_, err := services.GetQuestionBoxService().AnswerList(question, 0, 10)
	assert.Equal(t, err, nil)
	_, err = services.GetQuestionBoxService().AnswerList(question, -1, 10)
	assert.NotEqual(t, err, nil)
	_, err = services.GetQuestionBoxService().AnswerList(question, -1, -10)
	assert.NotEqual(t, err, nil)
	_, err = services.GetQuestionBoxService().AnswerList(question, 1, -10)
	assert.NotEqual(t, err, nil)
	_, err = services.GetQuestionBoxService().AnswerList(nil, 1, 10)
	assert.NotEqual(t, err, nil)

}

func TestMyAnswerList(t *testing.T) {
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)
	_, err = services.GetQuestionBoxService().MyAnswerList(&user, 0, 10)
	assert.Equal(t, err, nil)
	_, err = services.GetQuestionBoxService().MyAnswerList(&user, -1, 10)
	assert.NotEqual(t, err, nil)
	_, err = services.GetQuestionBoxService().MyAnswerList(&user, -1, -10)
	assert.NotEqual(t, err, nil)
	_, err = services.GetQuestionBoxService().MyAnswerList(&user, 1, -10)
	assert.NotEqual(t, err, nil)
	_, err = services.GetQuestionBoxService().MyAnswerList(nil, 1, -10)
	assert.NotEqual(t, err, nil)
}

func BenchmarkMyAnswerList(b *testing.B) {
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(b, err, nil)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := services.GetQuestionBoxService().MyAnswerList(&user, int64(rand.Int()), int64(rand.Int()))
			assert.Equal(b, err, nil)
		}
	})

}

func BenchmarkAnswerList(b *testing.B) {
	questionId, _ := primitive.ObjectIDFromHex("645df98ee933a87982169266")
	question := &models.QuestionBoxQuestion{
		ID: questionId,
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := services.GetQuestionBoxService().AnswerList(question, int64(rand.Int()), int64(rand.Int()))
			assert.Equal(b, err, nil)
		}
	})

}

func TestUpdateAnswerContent(t *testing.T) {
	question, err := getOneQuestion("TestUpdateAnswerContent")
	assert.Equal(t, err, nil)

	questionID, err := services.GetQuestionBoxService().NewQuestion(question)
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
	}()
	assert.Equal(t, err, nil)

	answer, err := getOneAnswer(question)
	assert.Equal(t, err, nil)

	answerID, err := services.GetQuestionBoxService().NewAnswer(answer)
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestionBoxAnswerByID(answerID)
	}()
	assert.Equal(t, err, nil)

	answer.Content = "Update 1"

	err = services.GetQuestionBoxService().UpdateAnswerContent(answer)
	assert.Equal(t, err, nil)

	updatedAnswer, err := services.GetQuestionBoxService().QueryAnswerByID(answerID)
	assert.Equal(t, err, nil)
	assert.Equal(t, updatedAnswer.Content, "Update 1")

	answer.Content = ""
	err = services.GetQuestionBoxService().UpdateAnswerContent(answer)
	assert.Equal(t, err, errors.New("更新的回答为空"))
}

func BenchmarkUpdateAnswerContent(b *testing.B) {
	question, err := getOneQuestion("BenchmarkUpdateAnswerContent")
	assert.Equal(b, err, nil)

	questionID, err := services.GetQuestionBoxService().NewQuestion(question)
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestion(questionID)
	}()
	assert.Equal(b, err, nil)

	answer, err := getOneAnswer(question)
	assert.Equal(b, err, nil)

	answerID, err := services.GetQuestionBoxService().NewAnswer(answer)
	defer func() {
		_ = services.GetQuestionBoxService().DeleteQuestionBoxAnswerByID(answerID)
	}()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			err = services.GetQuestionBoxService().UpdateAnswerContent(answer)
			assert.Equal(b, err, nil)
		}
	})
}

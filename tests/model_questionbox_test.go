package tests

import (
	"fmt"
	"melodie-site/server/models"
	"melodie-site/server/services"

	"testing"
	"time"

	"github.com/go-playground/assert/v2"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewQuestionboxAnswer(t *testing.T) {
	admin, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)
	Questionid, _ := primitive.ObjectIDFromHex("0")

	answer := &models.QuestionBoxAnswer{
		UserID:     admin.ID,
		CreateTime: uint64(time.Now().UnixMilli()),
		Content:    "华东师范大学中山北路校区在上海市普陀区中山北路，学校很小，方便赶早八",
		QuestionID: Questionid,

		Respondant: models.PersonalInfo{
			CEEPlace:  "河南省",
			Subject:   "理科",
			Age:       20,
			Gender:    "女",
			Situation: "参加2020年高考",
		},
	}

	err = services.GetQuestionBoxService().NewAnswer(answer)

	assert.Equal(t, err, nil)
}
func TestQueryAnswerByID(t *testing.T) {
	answerID, err := primitive.ObjectIDFromHex("64561d426097e10ebeda2b22")
	assert.Equal(t, err, nil)
	answer, err := services.GetQuestionBoxService().QueryAnswerByID(answerID)
	assert.Equal(t, err, nil)
	fmt.Print(answer.Content)
}
func TestGetAnswerList(t *testing.T) {
	Questionid, _ := primitive.ObjectIDFromHex("000000000000000000000000")
	question := &models.QuestionBoxQuestion{
		ID:          Questionid,
		Title:       "关于学校怎么样",
		Description: "学长学姐的精神状态怎么样",
	}
	answers, err := services.GetQuestionBoxService().AnswerList(question, 0)
	if answers==nil{
		fmt.Print("没有答案")
	}
	assert.Equal(t, err, nil)
	for _, answer := range answers {
		fmt.Println(answer.Content)
	}
}

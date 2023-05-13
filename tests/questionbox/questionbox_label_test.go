package tests

import (
	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"testing"
)

func TestNewLabel(t *testing.T) {
	labels := make([]*models.QuestionLabel, 3)
	for index := range labels {
		labels[index] = new(models.QuestionLabel)
		labels[index].Content = "Hello World"
		labels[index].UserID = primitive.NewObjectID()
		id, _ := primitive.ObjectIDFromHex("6453647c6fe2811ed76a9c04")
		labels[index].Questions = []models.QuestionInLabelInfo{
			{ID: id, Name: "My Question", HasRead: false},
		}
	}

	labelIDs, err := services.GetQuestionBoxService().NewLabels(labels)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, labelIDs, nil)
	println(labelIDs)
}

func TestUpdateLabelContent(t *testing.T) {
	labels, err := services.GetQuestionBoxService().QueryLabelsFromUser(
		&models.User{ID: primitive.NilObjectID},
		0, 0)
	assert.Equal(t, err, nil)

	label := labels[0]

	// 如果要更新内容
	label.Content = "Hello From World"

	err = services.GetQuestionBoxService().UpdateLabelContent(label)
	assert.Equal(t, err, nil)
}

func TestAddQuestionInLabel(t *testing.T) {
	questionID, _ := primitive.ObjectIDFromHex("64536b800a3da36ef0a12770")
	question, err := services.GetQuestionBoxService().QueryQuestionByID(questionID)
	assert.Equal(t, err, nil)

	questionInfo := &models.QuestionInLabelInfo{
		ID:      question.ID,
		Name:    question.Title,
		HasRead: false,
	}

	labelID, _ := primitive.ObjectIDFromHex("645fec942a221656eac2680b")

	err = services.GetQuestionBoxService().AddQuestionInLabel(labelID, questionInfo)
	assert.Equal(t, err, nil)

	label, err := services.GetQuestionBoxService().QueryLabelByID(labelID)
	assert.Equal(t, err, nil)
	assert.Equal(t, label.Content, "Hello From World")
}

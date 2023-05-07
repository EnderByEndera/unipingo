package routers

import (
	"encoding/json"
	"errors"

	"melodie-site/server/models"
	"melodie-site/server/services"
	"melodie-site/server/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewQuestionBoxAnswer(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// TODO 这里是不是要写一个questionboxanswerRequest结构体
	req := &models.QuestionBoxAnswer{}
	err = json.Unmarshal(dataBytes, req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	ans := &models.QuestionBoxAnswer{
		UserID:     userID,
		Content:    req.Content,
		QuestionID: req.QuestionID,
		School:     req.School,
		Major:      req.Major,
		Respondant: models.PersonalInfo{
			CEEPlace:  req.Respondant.CEEPlace,
			Subject:   req.Respondant.Subject,
			Age:       req.Respondant.Age,
			Gender:    req.Respondant.Gender,
			Situation: req.Respondant.Situation,
		},
	}
	err = services.GetQuestionBoxService().NewAnswer(ans)

	if err != nil {
		c.JSON(http.StatusBadRequest, makeResponse(false, err, nil))
		return
	}
}

func GetAnswerList(c *gin.Context) {

	// TODO  前端都传递哪些参数呢（目前用到questionID、page）
	questionboxID, err := primitive.ObjectIDFromHex(c.Query("questionboxID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, makeResponse(false, err, nil))
		return
	}
	question, err := services.GetQuestionBoxService().QueryQuestionByID(questionboxID)
	if err != nil {
		c.JSON(http.StatusBadRequest, makeResponse(false, errors.New("问题不存在"), nil))
		return

	}
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, makeResponse(false, err, nil))
		return
	}

	answers, err := services.GetQuestionBoxService().AnswerList(question, int64(page))
	if err != nil {
		c.JSON(http.StatusNotFound, makeResponse(false, errors.New("数据库查询错误"), nil))
	} else {
		c.JSON(http.StatusOK, makeResponse(true, nil, answers))
	}

}

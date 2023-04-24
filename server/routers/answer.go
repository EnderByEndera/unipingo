package routers

import (
	"encoding/json"
	"errors"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"melodie-site/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTopics(c *gin.Context) {
	topics := services.GetAnswersService().GetAllTopics()
	c.JSON(http.StatusOK, makeResponse(true, nil, topics))
}

func NewAnswer(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	reqStruct := &models.NewAnswerRequest{}
	err = json.Unmarshal(dataBytes, reqStruct)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	var entityName string
	if reqStruct.Category == models.AnswerAboutHEI {
		entityName, err = services.GetHEIService().GetHEIName(reqStruct.EntityID)
	} else {
		entityName, err = services.GetHEIService().GetHEIName(reqStruct.EntityID)
	}
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ans := &models.Answer{
		UserID:   userID,
		Content:  reqStruct.Content,
		Question: reqStruct.Question,
		Category: reqStruct.Category,
		BelongsTo: models.EntityWithName{
			ID:   reqStruct.EntityID,
			Name: entityName,
		},
	}
	id, err := services.GetAnswersService().NewAnswer(ans)
	if err != nil {
		c.JSON(http.StatusOK, makeResponse(true, nil, id))
		return
	}
}

func GetAnswersRelatedToHEIOrMajor(c *gin.Context) {
	entityIDStr := c.Query("entityID")
	category := c.Query("category")
	question := c.Query("question")
	if entityIDStr == "" || category == "" {
		c.JSON(http.StatusBadRequest, makeResponse(false, errors.New("未选择id或请求类型"), nil))
		return
	}
	if question == "" {
		c.JSON(http.StatusBadRequest, makeResponse(false, errors.New("未指定问题"), nil))
		return
	}
	entityID, err := primitive.ObjectIDFromHex(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, makeResponse(false, errors.New("entityID解码失败"), nil))
		return
	}
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, makeResponse(false, errors.New("userID解码失败"), nil))
		return
	}
	answers, err := services.GetAnswersService().GetAnswersRelatedToHEIOrMajor(entityID, models.AnswerCategory(category), question)
	ansResp := &models.AnswersResponse{Answers: answers}
	ansResp.Init()
	// ansResp.Approved[]
	// for _, ans := range answers {
	// 	ans.ApprovedUsers
	// }
	utils.ForEach(answers, func(i int, ans *models.Answer) {
		if utils.InArr(ans.ApprovedUsers, userID) {
			ansResp.Approved[ans.ID] = 1
		}
		if utils.InArr(ans.DisapprovedUsers, userID) {
			ansResp.Disapproved[ans.ID] = 1
		}
		if utils.InArr(ans.FavoritedUsers, userID) {
			ansResp.Favorited[ans.ID] = 1
		}
	})
	if err != nil {
		c.JSON(http.StatusNotFound, makeResponse(false, err, nil))
	} else {
		c.JSON(http.StatusOK, makeResponse(true, nil, ansResp))
	}
}

func ApproveOrDisapproveAnswer(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	reqStruct := &models.ApproveOrDisapproveAnswerRequest{}
	err = json.Unmarshal(dataBytes, reqStruct)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	var ret models.StatusReport
	if reqStruct.Approve {
		ret = services.GetAnswersService().ApproveAnswer(userID, reqStruct.AnsID)
	} else {
		ret = services.GetAnswersService().DisApproveAnswer(userID, reqStruct.AnsID)
	}
	if ret.Error != nil {
		c.JSON(http.StatusBadRequest, makeResponse(ret.Error == nil, ret.Error, ret))
	} else {
		c.JSON(http.StatusOK, makeResponse(ret.Error == nil, ret.Error, ret))
	}

	return
}

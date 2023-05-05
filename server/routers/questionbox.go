package routers

import (
	"encoding/json"
	"errors"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"melodie-site/server/svcerror"
	"melodie-site/server/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewQuestion(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	newQuestionReq := new(models.NewQuestionReq)
	err = json.Unmarshal(data, newQuestionReq)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	// 校验学校或专业ID是否正确
	if newQuestionReq.School.ID != primitive.NilObjectID {
		hei, err := services.GetHEIService().GetHEI(newQuestionReq.School.ID)
		if err != nil {
			c.Error(svcerror.New(http.StatusBadRequest, err))
			return
		}
		if hei.Name != newQuestionReq.School.Name {
			err = errors.New("HEI名称和ID对应失败")
			c.Error(svcerror.New(http.StatusBadRequest, err))
			return
		}
	} else if newQuestionReq.Major.ID != primitive.NilObjectID {
		major, err := services.GetMajorService().GetMajor(newQuestionReq.Major.ID)
		if err != nil {
			c.Error(svcerror.New(http.StatusBadRequest, err))
			return
		}
		if major.Name != newQuestionReq.Major.Name {
			err = errors.New("专业名称和ID对应失败")
			c.Error(svcerror.New(http.StatusBadRequest, err))
			return
		}
	} else {
		// 学校和专业必须二选一，否则Request失败
		err = errors.New("学校和专业ID均为空")
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}
	// 校验userID是否有效
	user, err := services.GetAuthService().GetUserByID(userID)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	question := &models.QuestionBoxQuestion{
		UserID:      user.ID,
		Title:       newQuestionReq.Title,
		Description: newQuestionReq.Description,
		School:      newQuestionReq.School,
		Major:       newQuestionReq.Major,
		Questioner:  newQuestionReq.Questioner,
	}

	docID, err := services.GetQuestionBoxService().NewQuestion(question)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"docID": docID,
	})
}

func QueryQuestionByID(c *gin.Context) {
	questionID, err := primitive.ObjectIDFromHex(c.Query("question_id"))
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	_, err = services.GetAuthService().GetUserByID(userID)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	question, err := services.GetQuestionBoxService().QueryQuestionByID(questionID)
	if err != nil {
		c.Error(svcerror.New(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, question)
}

func QueryQuestionList(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	user, err := services.GetAuthService().GetUserByID(userID)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	pageNum, err := strconv.Atoi(c.Query("pageNum"))
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	questions, err := services.GetQuestionBoxService().QuestionList(user, int64(page), int64(pageNum))
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, questions)
}

func UpdateQuestionDescription(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	_, err = services.GetAuthService().GetUserByID(userID)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	question := new(models.QuestionBoxQuestion)
	if err = c.ShouldBind(question); err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	err = services.GetQuestionBoxService().UpdateQuestionDescription(question)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"update": true,
	})
}

func UpdateQuestionSchoolOrMajor(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	if _, err = services.GetAuthService().GetUserByID(userID); err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	question := new(models.QuestionBoxQuestion)
	if err = c.ShouldBind(question); err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	if err = services.GetQuestionBoxService().UpdateQuestionSchoolOrMajor(question); err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"update": true,
	})
}

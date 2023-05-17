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

// NewQuestion 新增问题接口
// @Summary 新增一个问题
// @Schemes
// @Description 通过提问表单，在数据库中新增一个问题
// @Tags questionbox
// @Param newQuestionReq body models.NewQuestionReq true "新增问题请求"
// @Accept application/json
// @Produce application/json
// @Success 200 {object} models.NewQuestionRes "新增问题响应“
// @Failure 400 {object} svcerror.SvcErr
// @Failure 500 {object} svcerror.SvcErr
// @Router /questionbox/question/new [post]
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

	c.JSON(http.StatusOK, models.NewQuestionRes{QuestionID: docID})
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

	questions, err := services.GetQuestionBoxService().QueryQuestionsFromUser(user, int64(page), int64(pageNum))
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

func NewQuestionBoxAnswer(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	req := &models.QuestionBoxAnswerReq{}
	err = json.Unmarshal(dataBytes, req)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	// 校验学校或专业ID是否正确
	if req.School.ID != primitive.NilObjectID {
		hei, err := services.GetHEIService().GetHEI(req.School.ID)
		if err != nil {
			c.Error(svcerror.New(http.StatusBadRequest, err))
			return
		}
		if hei.Name != req.School.Name {
			err = errors.New("HEI名称和ID对应失败")
			c.Error(svcerror.New(http.StatusBadRequest, err))
			return
		}
	} else if req.Major.ID != primitive.NilObjectID {
		major, err := services.GetMajorService().GetMajor(req.Major.ID)
		if err != nil {
			c.Error(svcerror.New(http.StatusBadRequest, err))
			return
		}
		if major.Name != req.Major.Name {
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
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}
	// 校验userID是否有效
	user, err := services.GetAuthService().GetUserByID(userID)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}
	ans := &models.QuestionBoxAnswer{
		UserID:     user.ID,
		Content:    req.Content,
		QuestionID: req.QuestionID,
		School:     req.School,
		Major:      req.Major,
		Respondant: req.Respondant,
	}
	docID, err := services.GetQuestionBoxService().NewAnswer(ans)

	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"docID": docID,
	})
}

func QueryAnswerByID(c *gin.Context) {
	answerID, err := primitive.ObjectIDFromHex(c.Query("answer_id"))
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

	question, err := services.GetQuestionBoxService().QueryQuestionByID(answerID)
	if err != nil {
		c.Error(svcerror.New(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, question)
}

func GetAnswerList(c *gin.Context) {
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

	questionboxID, err := primitive.ObjectIDFromHex(c.Query("questionboxID"))
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}
	question, err := services.GetQuestionBoxService().QueryQuestionByID(questionboxID)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return

	}
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}
	pageNums := c.Query("pageNum")
	pageNum, err := strconv.Atoi(pageNums)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	answers, err := services.GetQuestionBoxService().AnswerList(question, int64(page), int64(pageNum))
	if err != nil {
		c.JSON(http.StatusNotFound, makeResponse(false, errors.New("数据库查询错误"), nil))
	} else {
		c.JSON(http.StatusOK, makeResponse(true, nil, answers))
	}
}

func GetMyAnswerList(c *gin.Context) {
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

	answers, err := services.GetQuestionBoxService().MyAnswerList(user, int64(page), int64(pageNum))
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, answers)
}

func UpdateAnswerContent(c *gin.Context) {
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

	answer := &models.QuestionBoxAnswer{}
	if err = c.ShouldBind(answer); err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	err = services.GetQuestionBoxService().UpdateAnswerContent(answer)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"update": true,
	})

}

func NewLabels(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	user, err := services.GetAuthService().GetUserByID(userID)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	labels := make([]*models.QuestionLabel, 0)
	err = c.BindJSON(&labels)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	for index := range labels {
		labels[index].UserID = user.ID
	}

	labelIDs, err := services.GetQuestionBoxService().NewLabels(labels)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"labels": labelIDs,
	})
}

func GetLabelsFromUser(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	user, err := services.GetAuthService().GetUserByID(userID)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	req := new(models.GetLabelsFromUserRequest)
	err = c.BindJSON(req)

	labels, err := services.GetQuestionBoxService().QueryLabelsFromUser(user, req.Page, req.PageNum)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"labels": labels,
	})
}

func GetLabelFromQuestion(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	user, err := services.GetAuthService().GetUserByID(userID)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	type GetLabelFromQuestionReq struct {
		QuestionID primitive.ObjectID
		Page       int64
		PageNum    int64
	}

	req := new(GetLabelFromQuestionReq)
	err = c.BindJSON(req)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	question, err := services.GetQuestionBoxService().QueryQuestionByID(req.QuestionID)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
	}

	labels, err := services.GetQuestionBoxService().QueryLabelsFromQuestion(user, question, req.Page, req.PageNum)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
	}

	c.JSON(http.StatusOK, gin.H{
		"labels": labels,
	})
}

func DeleteLabel(c *gin.Context) {
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

	label := new(models.QuestionLabel)
	err = c.BindJSON(&label)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	err = services.GetQuestionBoxService().DeleteLabel(label.ID)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"delete": true,
	})
}

func UpdateLabel(c *gin.Context) {
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

	label := new(models.QuestionLabel)
	err = c.BindJSON(&label)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
	}

	err = services.GetQuestionBoxService().UpdateLabelContent(label)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
	}

	c.JSON(http.StatusOK, gin.H{
		"update": true,
	})
}

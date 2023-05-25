package routers

import (
	"encoding/json"
	"errors"
	"melodie-site/server/routers/types"

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
// @summary 新增一个问题
// @description 通过提问表单，在数据库中新增一个问题
// @tags questionbox,question
// @param x-access-token header string true "JWT Token"
// @param newQuestionReq body types.NewQuestionReq true "新增问题请求"
// @accept application/json
// @produce application/json
// @success 200 {object} types.NewQuestionRes "新增问题响应“
// @failure 400 {object} svcerror.SvcErr "请求格式问题"
// @failure 401 {object{ svcerror.SvcErr "用户认证失败"
// @failure 500 {object} svcerror.SvcErr "服务器内部问题"
// @router /questionbox/question/new [post]
func NewQuestion(c *gin.Context) {
	newQuestionReq := new(types.NewQuestionReq)
	err := c.ShouldBindJSON(newQuestionReq)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	// 校验学校或专业ID是否正确
	if newQuestionReq.School.ID != primitive.NilObjectID {
		hei, err := services.GetHEIService().GetHEI(newQuestionReq.School.ID)
		if err != nil {
			c.Error(svcerror.New(http.StatusBadRequest, errors.New("未知学校ID")))
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
			c.Error(svcerror.New(http.StatusBadRequest, errors.New("未知专业ID")))
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
		c.Error(svcerror.New(http.StatusUnauthorized, err))
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
		AskTo:       newQuestionReq.AskTo,
		AskTags:     newQuestionReq.AskTags,
	}

	docID, err := services.GetQuestionBoxService().NewQuestion(question)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, types.NewQuestionRes{QuestionID: docID})
}

// QueryQuestionByID 问题ID查询接口
// @summary 问题ID查询接口
// @description 根据ID查询一个问题
// @tags questionbox,question
// @param x-access-token header string true "JWT Token"
// @param question_id query string true "问题ID"
// @produce application/json
// @success 200 {object} types.QueryQuestionByIDRes "ID对应问题响应"
// @failure 400 {object} svcerror.SvcErr "请求格式问题"
// @failure 401 {object{ svcerror.SvcErr "用户认证失败"
// @failure 500 {object} svcerror.SvcErr "服务器内部问题"
// @router /questionbox/question/query [get]
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

	c.JSON(http.StatusOK, types.QueryQuestionByIDRes{
		Question: question,
	})
}

// QueryMyQuestionList 用户问题列表查询接口
// @summary 用户问题列表查询接口
// @description 查询某用户的所有问题
// @tags questionbox,question
// @param x-access-token header string true "JWT Token"
// @param page query uint64 true "查询页数"
// @param pageNum query uint64 true "一页需要查询的问题数量"
// @produce application/json
// @success 200 {object} types.QueryQuestionListRes "用户对应问题响应"
// @failure 400 {object} svcerror.SvcErr "请求格式问题"
// @failure 401 {object{ svcerror.SvcErr "用户认证失败"
// @failure 500 {object} svcerror.SvcErr "服务器内部问题"
// @router /questionbox/question/list [get]
func QueryMyQuestionList(c *gin.Context) {
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

	c.JSON(http.StatusOK, types.QueryQuestionListRes{
		Questions: questions,
		NextPage:  page + 1,
	})
}

// UpdateQuestionDescription 更新问题描述接口
// @summary 更新问题描述接口
// @description 更新某问题的描述信息
// @tags questionbox,question
// @param x-access-token header string true "JWT Token"
// @param updateQuestionDescriptionReq body types.UpdateQuestionDescriptionReq true "更新问题描述请求"
// @produce application/json
// @success 200 {object} types.UpdateQuestionDescriptionRes "更新问题描述响应"
// @failure 400 {object} svcerror.SvcErr "请求格式问题"
// @failure 401 {object{ svcerror.SvcErr "用户认证失败"
// @failure 500 {object} svcerror.SvcErr "服务器内部问题"
// @router /questionbox/question/description/update [post]
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

	questionReq := new(types.UpdateQuestionDescriptionReq)
	if err = c.ShouldBindJSON(questionReq); err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	err = services.GetQuestionBoxService().UpdateQuestionDescription(questionReq.Question)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, types.UpdateQuestionDescriptionRes{
		Update: true,
	})
}

// UpdateQuestionSchoolOrMajor 更新问题学校或专业接口
// @summary 更新问题学校或专业接口
// @description 更新某问题的询问学校或专业
// @tags questionbox,question
// @param x-access-token header string true "JWT Token"
// @param updateQuestionSchoolOrMajorReq body types.UpdateQuestionSchoolOrMajorReq true "更新问题学校或专业请求"
// @produce application/json
// @success 200 {object} types.UpdateQuestionSchoolOrMajorRes "更新问题学校或专业响应"
// @failure 400 {object} svcerror.SvcErr "请求格式问题"
// @failure 401 {object{ svcerror.SvcErr "用户认证失败"
// @failure 500 {object} svcerror.SvcErr "服务器内部问题"
// @router /questionbox/question/school/update [post]
// @router /questionbox/question/major/update [post]
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

	req := new(types.UpdateQuestionSchoolOrMajorReq)
	if err = c.ShouldBind(req); err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	if err = services.GetQuestionBoxService().UpdateQuestionSchoolOrMajor(req.Question); err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, types.UpdateQuestionSchoolOrMajorRes{
		Update: true,
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

	// 校验学校或专业ID是否存在且正确
	if req.School.ID == primitive.NilObjectID {
		err = errors.New("HEI的ID为空")
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}
	if req.Major.ID == primitive.NilObjectID {
		err = errors.New("major的ID为空")
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

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

	answer, err := services.GetQuestionBoxService().QueryAnswerByID(answerID)
	if err != nil {
		c.Error(svcerror.New(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, answer)
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

	answer := new(models.QuestionBoxAnswer)
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

// ReadAnswerByUser 用户已读某一回答接口
// @summary 用户已读某一回答接口
// @description 当用户读取某个回答的页面时调用该接口
// @tags questionbox,answer
// @param x-access-token header string true "JWT Token"
// @param newLabelsReq body types.ReadAnswerByUserReq true "用户读取某回答请求，需包含回答ID"
// @produce application/json
// @success 200 {object} types.ReadAnswerByUserRes "用户读取某回答响应"
// @failure 400 {object} svcerror.SvcErr "请求格式问题"
// @failure 401 {object{ svcerror.SvcErr "用户认证失败"
// @failure 500 {object} svcerror.SvcErr "服务器内部问题"
// @router /questionbox/label/new [post]
func ReadAnswerByUser(c *gin.Context) {
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

	req := new(types.ReadAnswerByUserReq)
	err = c.ShouldBindJSON(req)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	answer, err := services.GetQuestionBoxService().QueryAnswerByID(req.AnswerID)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	err = services.GetQuestionBoxService().ReadAnswerByUser(user.ID, answer)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, types.ReadAnswerByUserRes{
		Log: true,
	})
}

// NewLabels 新建问题标签（文件夹）接口
// @summary 新建问题标签（文件夹）接口
// @description 新建多个问题标签（文件夹）
// @tags questionbox,label
// @param x-access-token header string true "JWT Token"
// @param newLabelsReq body types.NewLabelsReq true "新建问题标签（文件夹）请求"
// @produce application/json
// @success 200 {object} types.NewLabelsRes "新建问题标签（文件夹）响应"
// @failure 400 {object} svcerror.SvcErr "请求格式问题"
// @failure 401 {object{ svcerror.SvcErr "用户认证失败"
// @failure 500 {object} svcerror.SvcErr "服务器内部问题"
// @router /questionbox/label/new [post]
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

	req := new(types.NewLabelsReq)
	err = c.ShouldBindJSON(req)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	for index := range req.Labels {
		req.Labels[index].UserID = user.ID
	}

	labelIDs, err := services.GetQuestionBoxService().NewLabels(req.Labels)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, types.NewLabelsRes{
		LabelIDs: labelIDs,
	})
}

// GetLabelsFromUser 根据用户获取标签（文件夹）接口
// @summary 根据用户获取标签（文件夹）接口
// @description 根据当前用户获取当前用户的所有标签（文件夹）
// @tags questionbox,label
// @param x-access-token header string true "JWT Token"
// @param page query uint64 true "标签（文件夹）页码"
// @param pageNum query uint64 true "标签（文件夹）每页包含个数"
// @produce application/json
// @success 200 {object} types.GetLabelsFromUserRes "获取用户对应标签（文件夹）响应"
// @failure 400 {object} svcerror.SvcErr "请求格式问题"
// @failure 401 {object{ svcerror.SvcErr "用户认证失败"
// @failure 500 {object} svcerror.SvcErr "服务器内部问题"
// @router /questionbox/label/user/get [get]
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

	labels, err := services.GetQuestionBoxService().QueryLabelsFromUser(user, int64(page), int64(pageNum))
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, types.GetLabelsFromUserRes{
		Labels: labels,
	})
}

// GetLabelsFromQuestion 根据问题获取标签（文件夹）接口
// @summary 根据问题获取标签（文件夹）接口
// @description 根据当前问题获取对应的所有标签（文件夹）
// @tags questionbox,label
// @param x-access-token header string true "JWT Token"
// @param page query uint64 true "标签（文件夹）页码"
// @param pageNum query uint64 true "标签（文件夹）每页包含个数"
// @param getLabelsFromQuestionReq body types.GetLabelsFromQuestionReq true "所搜寻问题ID"
// @produce application/json
// @success 200 {object} types.GetLabelsFromQuestionRes "获取问题对应标签（文件夹）响应"
// @failure 400 {object} svcerror.SvcErr "请求格式问题"
// @failure 401 {object{ svcerror.SvcErr "用户认证失败"
// @failure 500 {object} svcerror.SvcErr "服务器内部问题"
// @router /questionbox/label/question/get [post]
func GetLabelsFromQuestion(c *gin.Context) {
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

	req := new(types.GetLabelsFromQuestionReq)
	err = c.ShouldBindJSON(req)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	question, err := services.GetQuestionBoxService().QueryQuestionByID(req.QuestionID)
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

	labels, err := services.GetQuestionBoxService().QueryLabelsFromQuestion(user, question, int64(page), int64(pageNum))
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, types.GetLabelsFromQuestionRes{
		Labels: labels,
	})
}

// DeleteLabel 删除标签（文件夹）接口
// @summary 删除标签（文件夹）接口
// @description 根据标签（文件夹）ID删除对应的标签（文件夹）
// @tags questionbox,label
// @param x-access-token header string true "JWT Token"
// @param page query string true "标签（文件夹）ID"
// @produce application/json
// @success 200 {object} types.DeleteLabelRes "获取问题对应标签（文件夹）响应"
// @failure 400 {object} svcerror.SvcErr "请求格式问题"
// @failure 401 {object{ svcerror.SvcErr "用户认证失败"
// @failure 500 {object} svcerror.SvcErr "服务器内部问题"
// @router /questionbox/label/question/delete [post]
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

	labelID := primitive.ObjectID{}
	err = c.ShouldBindJSON(&labelID)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	err = services.GetQuestionBoxService().DeleteLabel(labelID)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, types.DeleteLabelRes{
		LabelID: labelID,
	})
}

// UpdateLabelContent 更新标签（文件夹）接口
// @summary 更新标签（文件夹）接口
// @description 更新标签（文件夹）的标题（内容）
// @tags questionbox,label
// @param x-access-token header string true "JWT Token"
// @param page query string true "标签（文件夹）ID"
// @param updateLabelContentReq body types.UpdateLabelContentReq true "标签（文件夹）更新内容请求"
// @produce application/json
// @success 200 {object} types.UpdateLabelContentRes "标签（文件夹）更新响应"
// @failure 400 {object} svcerror.SvcErr "请求格式问题"
// @failure 401 {object{ svcerror.SvcErr "用户认证失败"
// @failure 500 {object} svcerror.SvcErr "服务器内部问题"
// @router /questionbox/label/content/update [post]
func UpdateLabelContent(c *gin.Context) {
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

	req := new(types.UpdateLabelContentReq)
	err = c.ShouldBindJSON(req)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	label := &models.QuestionBoxLabel{
		ID:      req.LabelID,
		Content: req.Content,
	}

	err = services.GetQuestionBoxService().UpdateLabelContent(label)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, types.UpdateLabelContentRes{
		LabelID: label.ID,
	})
}

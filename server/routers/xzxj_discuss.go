package routers

import (
	"github.com/gin-gonic/gin"
	"melodie-site/server/models"
	"melodie-site/server/routers/types"
	"melodie-site/server/services"
	"melodie-site/server/svcerror"
	"melodie-site/server/utils"
	"net/http"
)

// AddOrUpdateXZXJUser 添加学长学姐说用户接口
// @summary 添加学长学姐说用户接口
// @description 增加以为学长学姐说用户
// @tags xzxj_discuss
// @param x-access-token header string true "JWT Token"
// @param req body types.AddOrUpdateXZXJUserReq true "学长学姐说添加用户接口请求"
// @produce application/json
// @success 200 {object} types.AddOrUpdateXZXJUserRes "学长学姐说添加用户接口响应"
// @failure 400 {object} svcerror.SvcErr "请求格式问题"
// @failure 401 {object{ svcerror.SvcErr "用户认证失败"
// @failure 500 {object} svcerror.SvcErr "服务器内部问题"
// @router /xzxjdiscuss/add [post]
// @router /xzxjdiscuss/update [post]
func AddOrUpdateXZXJUser(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	req := new(types.AddOrUpdateXZXJUserReq)
	err = c.ShouldBindJSON(req)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	xzxjUserForm := &models.XZXJUserFormMeta{
		XZXJUser: models.XZXJUser{
			UserID:            userID,
			Sections:          req.Sections,
			Picture:           req.Picture,
			Motto:             req.Motto,
			ManagedActivities: req.ManagedActivities,
			Experience:        req.Experience,
		},
		RealName: req.RealName,
		UserTags: req.UserTags,
	}

	userID, err = services.GetXZXJDiscussService().AddOrUpdateXZXJUser(xzxjUserForm)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, types.AddOrUpdateXZXJUserRes{
		UserID: userID,
	})
}

// QueryXZXJUserByUserID 查询学长学姐说用户接口
// @summary 查询学长学姐说用户接口
// @description 查询学长学姐说用户
// @tags xzxj_discuss
// @param x-access-token header string true "JWT Token"
// @produce application/json
// @success 200 {object} types.QueryXZXJUserByUserIDRes "查询学长学姐说用户接口响应"
// @failure 400 {object} svcerror.SvcErr "请求格式问题"
// @failure 401 {object{ svcerror.SvcErr "用户认证失败"
// @failure 500 {object} svcerror.SvcErr "服务器内部问题"
// @router /xzxjdiscuss/get [get]
func QueryXZXJUserByUserID(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	req := new(types.QueryXZXJUserByUserIDReq)
	err = c.ShouldBindJSON(req)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	xzxjUser, err := services.GetXZXJDiscussService().QueryXZXJUserByUserID(userID)

	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, types.QueryXZXJUserByUserIDRes{
		XZXJUser: xzxjUser,
	})
}

// DeleteXZXJUser 删除学长学姐说用户接口
// @summary 删除学长学姐说用户接口
// @description 删除学长学姐说用户
// @tags xzxj_discuss
// @param x-access-token header string true "JWT Token"
// @produce application/json
// @success 200 {object} types.DeleteXZXJUserRes "删除学长学姐说用户接口响应"
// @failure 400 {object} svcerror.SvcErr "请求格式问题"
// @failure 401 {object{ svcerror.SvcErr "用户认证失败"
// @failure 500 {object} svcerror.SvcErr "服务器内部问题"
// @router /xzxjdiscuss/delete [get]
func DeleteXZXJUser(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	err = services.GetXZXJDiscussService().DeleteXZXJUserByID(userID)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, types.DeleteXZXJUserRes{
		Deleted: true,
	})
}

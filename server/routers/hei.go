package routers

import (
	"encoding/json"
	"errors"
	"fmt"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"melodie-site/server/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetHEI(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, makeResponse(false, errors.New("学校id不能为空"), nil))
		return
	}
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, makeResponse(false, fmt.Errorf("无效学校id %s", idStr), nil))
		return
	}
	hei, err := services.GetHEIService().GetHEI(id)
	if err != nil {
		c.JSON(http.StatusNotFound, makeResponse(false, fmt.Errorf("未找到id为'%s'的学校", idStr), nil))
	} else {
		c.JSON(http.StatusOK, makeResponse(true, nil, hei))
	}

}

// 通过名称获取高等教育机构
// GET
// URL参数: ?name=<学校名称>
func GetHEIByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, makeResponse(false, errors.New("学校名称不能为空"), nil))
		return
	}
	hei, err := services.GetHEIService().GetHEIByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, makeResponse(false, errors.New(fmt.Sprintf("未找到名为'%s'的学校", name)), nil))
	} else {
		c.JSON(http.StatusOK, makeResponse(true, nil, hei))
	}

}

// 过滤高等教育机构
// 返回高等教育机构列表，以json形式
// url参数名称参照services.GetHEIService().FilterHEI()的输入参数
// 要求自测通过，对于异常输入，返回中文的报错提示。
func FilterHEI(c *gin.Context) {
	//拿到参数provincialLocation
	provincialLocation := c.Query("provincialLocation")

	//参数level
	s := c.Query("level")
	var levels models.HEILevel
	if s != "" {
		level, err := strconv.Atoi(s)
		//如果输入level不为数字，就出错
		if err != nil {
			c.JSON(http.StatusBadRequest, makeResponse(false, errors.New(fmt.Sprintf("level参数格式错误")), nil))
			return
		}
		levels = models.HEILevel(level)
	} else {
		levels = -1
	}
	//参数mode
	p := c.Query("mode")
	var modes models.HEIMode
	if p != "" {
		mode, err := strconv.Atoi(p)
		//如果输入mode不为数字，就出错
		if err != nil {
			c.JSON(http.StatusBadRequest, makeResponse(false, errors.New(fmt.Sprintf("mode参数格式错误")), nil))
			return
		}
		modes = models.HEIMode(mode)
	} else {
		modes = -1
	}
	//参数policy
	policy := c.Query("policy")

	//分页参数，当前是多少页。
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, makeResponse(false, err, nil))
		return
	}

	heis, err := services.GetHEIService().FilterHEI(provincialLocation, levels, modes, policy, int64(page))

	if err != nil {
		c.JSON(http.StatusNotFound, makeResponse(false, errors.New(fmt.Sprintf("数据库查询错误")), nil))
	} else {
		c.JSON(http.StatusOK, makeResponse(true, nil, heis))
	}
}

func AddHEIToCollection(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(500, err)
	}
	reqStruct := &models.EntityWithName{}
	err = json.Unmarshal(dataBytes, reqStruct)
	if err != nil {
		c.AbortWithError(500, err)
	}
	fmt.Printf("%+v\n", reqStruct)
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, makeResponse(false, err, nil))
		return
	}
	if ok, _ := services.GetAuthService().IsHEIOrMajorInCollection(userID, reqStruct.ID, models.CollectionItemHEI); ok {
		c.JSON(http.StatusConflict, makeResponse(false, errors.New("已经收藏过"), nil))
		return
	}
	err = services.GetAuthService().AddHEIOrMajorToCollection(
		userID,
		reqStruct.ID,
		models.CollectionItemHEI,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, makeResponse(false, err, nil))
		return
	}
	c.JSON(http.StatusOK, makeResponse(true, nil, "添加成功"))

}

func RemoveHEIFromCollection(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(500, err)
	}
	reqStruct := &models.EntityWithName{}
	err = json.Unmarshal(dataBytes, reqStruct)
	if err != nil {
		c.AbortWithError(500, err)
	}
	fmt.Printf("%+v\n", reqStruct)
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, makeResponse(false, err, nil))
		return
	}
	if ok, _ := services.GetAuthService().IsHEIOrMajorInCollection(userID, reqStruct.ID, models.CollectionItemHEI); !ok {
		c.JSON(http.StatusNotFound, makeResponse(false, errors.New("未被收藏，无法取消赞"), nil))
		return
	}
	err = services.GetAuthService().RemoveHEIOrMajorFromCollection(
		userID,
		reqStruct.ID,
		models.CollectionItemHEI,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, makeResponse(false, err, nil))
		return
	}
	c.JSON(http.StatusOK, makeResponse(true, nil, "取消收藏成功"))

}

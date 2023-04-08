package routers

import (
	"errors"
	"fmt"
	"melodie-site/server/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	// services.GetHEIService().FilterHEI()
}

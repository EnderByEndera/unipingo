package routers

import (
	"errors"
	"fmt"
	"melodie-site/server/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

//通过名称获得专业
// GET
// URL参数: ?name=<专业名称>
func GetMajorByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, makeResponse(false, errors.New("专业名称不能为空"), nil))
		return
	}
	major, err := services.GetMajorService().GetMajorByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, makeResponse(false, errors.New(fmt.Sprintf("未找到名为'%s'的专业", name)), nil))
	} else {
		c.JSON(http.StatusOK, makeResponse(true, nil, major))
	}
}

// 过滤major
// 返回major列表，以json形式
// url参数名称参照services.GetMajorService().FilterMajor()的输入参数
// 要求自测通过，对于异常输入，返回中文的报错提示。
func FilterMajor(c *gin.Context){
     category := c.Query("category")
	 majors, err := services.GetMajorService().FilterMajor(category)
	 if err != nil {
		c.JSON(http.StatusNotFound, makeResponse(false, errors.New(fmt.Sprintf("数据库查询错误")), nil))
	} else {
		c.JSON(http.StatusOK, makeResponse(true, nil, majors))
	}

}

package routers

import (
	"errors"
	"fmt"
	"melodie-site/server/services"
	"net/http"
	"strconv"
	"melodie-site/server/models"
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
	//拿到参数provincialLocation
	provincialLocation := c.Query("provincialLocation")

	//level和mode这两个参数要做类型转换，不知道学长有没有更简便的方法（我这里写的感觉笨笨的
	s := c.Query("level")
	var levels models.HEILevel
	if s != ""{
		level, err := strconv.Atoi(s)
		//如果输入level不为数字，就出错
		if err !=nil{
			c.JSON(http.StatusNotFound, makeResponse(false, errors.New(fmt.Sprintf("level参数格式错误")), nil))
			return
		}
		levels =models.HEILevel(level)
	}else{
		levels=-1
	}
	//参数mode
	p := c.Query("mode")
	var modes models.HEIMode
	if p!=""{
		mode, err := strconv.Atoi(p)
		//如果输入mode不为数字，就出错
		if err !=nil{
			c.JSON(http.StatusNotFound, makeResponse(false, errors.New(fmt.Sprintf("mode参数格式错误")), nil))
			return
		}
		modes =models.HEIMode(mode)
	}else{
		modes=-1
	}
    //参数policy
	policy := c.Query("policy")

	
	heis, err := services.GetHEIService().FilterHEI(provincialLocation, levels, modes, policy)
	//这里不太懂：如果service方法err不为空的话，都是什么情况？
	//(我这里测试：如果数据库中不存在相关数据(比如我设置level=3，这里没有符合条件的学校)，
	//services包里面FilterHEI方法中Find方法的返回值err为空，但是学长您写的GetHEIByName方法中如果没有对应学校会返回的err不为nil)
	if err !=nil{
		c.JSON(http.StatusNotFound, makeResponse(false, errors.New(fmt.Sprintf("数据库查询错误")), nil))
	}else{
		c.JSON(http.StatusOK, makeResponse(true, nil, heis))
	}
}

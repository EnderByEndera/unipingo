package routers

import (
	"encoding/json"
	"fmt"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTag(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(500, err)
	}
	tagReq := &models.Tag{}
	err = json.Unmarshal(dataBytes, tagReq)
	if err != nil {
		c.AbortWithError(500, err)
	}
	err = services.GetTagsService().CreateTag(tagReq)
	if err != nil {
		c.String(400, err.Error())
	} else {
		c.String(200, "Successfully created tag!")
	}
}

func GetTag(c *gin.Context) {
	tagID := c.Query("id")
	service := services.GetTagsService()
	id, err := strconv.Atoi(tagID)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	article, err := service.GetTag(id)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	c.JSON(200, article)
}

func GetAllTags(c *gin.Context) {
	tagsService := services.GetTagsService()
	tags, err := tagsService.GetAllTags()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, tags)
}

func UpdateTag(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(500, err)
	}
	req := &models.Tag{}
	err = json.Unmarshal(dataBytes, req)
	if err != nil {
		c.AbortWithError(500, err)
	}
	fmt.Println(req)
	services.GetTagsService().UpdateTag(req)
	c.String(200, "Successfully updated tag #%d!", req.ID)
}

package routers

import (
	"encoding/json"
	"fmt"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetArticle(c *gin.Context) {
	articleID := c.Query("id")
	articleService := services.GetArticleService()
	id, err := strconv.Atoi(articleID)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	article, err := articleService.GetArticle(id)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	c.JSON(200, article)
}

func GetAllArticles(c *gin.Context) {
	articleService := services.GetArticleService()
	articles, err := articleService.GetAllArticles()
	if err != nil {
		fmt.Println(err.Error())
		c.Error(err)
		return
	}
	c.JSON(200, articles)
}

func CreateArticle(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(500, err)
	}
	articleReq := &models.Article{}
	err = json.Unmarshal(dataBytes, articleReq)
	if err != nil {
		c.AbortWithError(500, err)
	}
	services.GetArticleService().CreateArticle(articleReq)
	c.String(200, "Successfully created article!")
}

func UpdateArticle(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(500, err)
	}
	articleReq := &models.Article{}
	err = json.Unmarshal(dataBytes, articleReq)
	if err != nil {
		c.AbortWithError(500, err)
	}
	services.GetArticleService().UpdateArticle(articleReq)
	c.String(200, "Successfully updated article!")
}

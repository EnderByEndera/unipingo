package routers

import (
	"melodie-site/server/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTopics(c *gin.Context) {
	topics := services.GetAnswersService().GetAllTopics()
	c.JSON(http.StatusOK, makeResponse(true, nil, topics))
}

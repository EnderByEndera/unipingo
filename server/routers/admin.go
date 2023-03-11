package routers

import (
	"encoding/json"
	"fmt"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"melodie-site/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewStuIDAuthProc(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	req := &models.NewStudentIdentityAuthenticationRequest{}
	err = json.Unmarshal(dataBytes, req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	auth, err := req.ToAuthStruct(claims.UserID)
	fmt.Println("auth started for user", auth.UserID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	code, err := services.GetAdminService().NewStuIDAuthProc(&auth)
	if err != nil {
		c.AbortWithError(code, err)
		return
	}
	c.String(http.StatusOK, "auth stream succeeded!")
}

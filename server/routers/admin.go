package routers

import (
	"encoding/json"
	"errors"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"melodie-site/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	code, err := services.GetAdminService().NewStuIDAuthProc(&auth)
	if err != nil {
		c.AbortWithError(code, err)
		return
	}
	services.GetAuthService().UpdateUserSchoolInfo(auth.UserID,
		&models.SchoolInfo{Name: req.SchoolName, Status: models.StudentIdentityPhotoUploaded})
	c.String(http.StatusOK, "auth stream succeeded!")
}

// 获取用户自己的学生身份认证
func GetStuIDAuthProc(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	proc, err := services.GetAdminService().GetStuIDAuthProc(userID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, proc)
}

func GetUnhandledProcs(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if claims.Role != models.RoleAdmin {
		c.AbortWithError(http.StatusForbidden, errors.New("您不是管理员，无法查看此信息"))
		return
	}
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	procs, err := services.GetAdminService().GetUnhandledStuIDAuthProcs()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, procs)
}

func SetStudentIDAuthStatus(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	req := &models.ModifyStuIDAuthStatRequest{}
	err = json.Unmarshal(dataBytes, req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = services.GetAdminService().UpdateStuIDAuthStatus(req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.String(http.StatusOK, "update succeeded!")
}

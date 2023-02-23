package routers

import (
	"fmt"
	"melodie-site/server/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendFile(ctx *gin.Context) {
	// splitted := strings.Split(ctx.Request.URL.Path, "/")
	fileName := ctx.Param("file")
	fmt.Println(fileName)
	if fileName == "" {
		ctx.String(http.StatusBadRequest, "file name empty")
		return
	}
	saveAsFile, err := services.GetFileFromOSS(fileName)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	} else {
		ctx.File(saveAsFile)
	}
}

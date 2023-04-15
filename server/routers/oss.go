package routers

import (
	"fmt"
	"melodie-site/server/services"
	"melodie-site/server/utils"
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

func UploadStaticFile(ctx *gin.Context) {
	f, file, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusOK, makeResponse(false, fmt.Errorf("file not in form"), nil))
		return
	}
	// file.Size/
	// fileName := ctx.Request.FormValue("fileName")
	err = utils.PutObject(
		utils.GetOSSHandler().Buckets.StaticFiles,
		file.Filename,
		file.Header.Get("Content-Type"),
		f,
		file.Size)
	if err != nil {
		ctx.JSON(http.StatusOK, makeResponse(false, err, nil))
		fmt.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, makeResponse(true, nil, nil))
	defer f.Close()
}

package services

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"melodie-site/server/utils"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func GetFileFromOSS(objectName string) (saveAsFile string, err error) {
	f, err := ioutil.TempFile(os.TempDir(), uuid.NewString())
	saveAsFile = f.Name()
	defer func() {
		f.Close()
	}()
	err = utils.GetOSSHandler().GetFileAndSave(objectName, saveAsFile)
	return
}

func GetStaticFileFromOSS(objectName string) (saveAsFile string, err error) {
	f, err := ioutil.TempFile(os.TempDir(), uuid.NewString())
	saveAsFile = f.Name()
	defer func() {
		f.Close()
	}()
	handler := utils.GetOSSHandler()
	err = handler.Client.FGetObject(context.Background(), handler.Buckets.StaticFiles, objectName,
		saveAsFile, minio.GetObjectOptions{})
	return
}

func UploadFileByHeaderToOSS(ctx *gin.Context, fileHeader *multipart.FileHeader) (fileName string, code int, err error) {
	f, err := ioutil.TempFile(os.TempDir(), uuid.NewString())
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	if err != nil {
		code = http.StatusInternalServerError
		return
	}
	err = ctx.SaveUploadedFile(fileHeader, f.Name())

	if err != nil {
		code = http.StatusInternalServerError
		return
	}
	fileName = uuid.NewString()
	utils.PutFile(utils.GetOSSHandler().Buckets.Files, fileName, f.Name(), fileHeader.Header.Get("Content-Type"))
	code = http.StatusOK
	return
}

func UploadMultipartFileToOSS(ext string, f multipart.File) (code int, err error) {
	buf := make([]byte, 512)

	length, err := f.Read(buf)
	if err != nil {
		fmt.Println(err)
		code = http.StatusBadRequest
		return
	}
	newOffset, err := f.Seek(0, 0)
	if err != nil || newOffset != 0 {
		err = errors.New(fmt.Sprintf("new offset is %v, other error: ", newOffset) + err.Error())
		code = http.StatusInternalServerError
		return
	}
	contentType := http.DetectContentType(buf[:length])
	if err != nil {
		code = http.StatusInternalServerError
		return
	}
	t0 := time.Now()
	ossHandler := utils.GetOSSHandler()

	fileObjectName := uuid.New().String() + "." + ext
	err = utils.PutObject(ossHandler.Buckets.Files, fileObjectName, contentType, f, -1)
	if err != nil {
		fmt.Println(err)
		code = http.StatusInternalServerError
		return
	}
	// t1 := time.Now()
	fmt.Println("t1", time.Since(t0))
	err = nil
	code = http.StatusOK
	return
}

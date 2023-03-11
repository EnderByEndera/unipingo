package services

import (
	"context"
	"errors"
	"melodie-site/server/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 管理员所用的一系列流程

type AdminService struct {
}

var adminService *AdminService

func GetAdminService() *AdminService {
	if adminService == nil {
		adminService = &AdminService{}
	}
	return adminService
}

func (adminService *AdminService) GetStuIDAuthProc(userID primitive.ObjectID) (authProc models.StudentIdentityAuthentication, err error) {
	filter := bson.M{"userID": userID}
	authProc = models.StudentIdentityAuthentication{}
	err = getCollection("auth").FindOne(context.TODO(), filter).Decode(&authProc)
	return
}

func (adminService *AdminService) NewStuIDAuthProc(auth *models.StudentIdentityAuthentication) (code int, err error) {
	_, err = adminService.GetStuIDAuthProc(auth.UserID)
	if err == nil {
		err = errors.New("record already exists for user: " + auth.UserID.Hex())
		code = http.StatusConflict // 返回资源冲突409。
		return
	}
	_, err = getCollection("auth").InsertOne(context.TODO(), auth)
	if err != nil {
		code = http.StatusBadRequest
		return
	}
	code = http.StatusOK
	return
}

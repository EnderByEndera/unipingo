package services

import (
	"context"
	"errors"
	"fmt"
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
		filter := bson.M{"userID": auth.UserID}
		fmt.Println(err, auth)
		err = getCollection("auth").FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": auth}).Err()

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

func (adminService *AdminService) GetUnhandledStuIDAuthProcs() (authProcs []models.StudentIdentityAuthentication, err error) {
	authProcs = make([]models.StudentIdentityAuthentication, 0)

	filter := bson.M{"status": models.StudentIdentityPhotoUploaded}
	cursor, err := getCollection("auth").Find(context.TODO(), filter)
	if err != nil {
		return
	}

	if err = cursor.All(context.TODO(), &authProcs); err != nil {
		panic(err)
	}
	return
}

func (adminService *AdminService) UpdateStuIDAuthStatus(req *models.ModifyStuIDAuthStatRequest) (err error) {
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return
	}
	proc, err := adminService.GetStuIDAuthProc(userID)
	if err == nil {
		filter := bson.M{"userID": userID}
		err = getCollection("auth").FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": bson.M{"status": req.Status, "suggestions": req.Suggestions}}).Err()
		if req.Status == models.StudentIdentityAuthenticated {
			err = GetAuthService().UpdateUserSchoolInfo(userID, &models.SchoolInfo{
				Name:     proc.SchoolName,
				Verified: true,
			})
			if err != nil {
				return
			}
		}
		return
	}
	return
}

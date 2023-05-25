package services

import (
	"context"
	"fmt"
	"melodie-site/server/db"
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

func (adminService *AdminService) GetStuIDAuthProc(id primitive.ObjectID) (authProc *models.StudentIdentityAuthentication, err error) {
	filter := bson.M{"_id": id}
	authProc = &models.StudentIdentityAuthentication{}
	res := db.GetCollection("auth").FindOne(context.TODO(), filter)
	err = res.Err()
	res.Decode(authProc)

	return
}

func (adminService *AdminService) NewStuIDAuthProc(auth *models.StudentIdentityAuthentication) (code int, err error) {
	_, err = db.GetCollection("auth").InsertOne(context.TODO(), auth)
	if err != nil {
		code = http.StatusBadRequest
		return
	}
	code = http.StatusOK
	return
}

func (adminService *AdminService) UpdateStuIDAuthProc(auth *models.StudentIdentityAuthentication) (code int, err error) {
	filter := bson.M{"_id": auth.ID}
	err = db.GetCollection("auth").FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": auth}).Err()
	code = http.StatusOK
	if err != nil {
		code = http.StatusNotFound
	}
	return
}

func (adminService *AdminService) GetUnhandledStuIDAuthProcs() (authProcs []models.StudentIdentityAuthentication, err error) {
	authProcs = make([]models.StudentIdentityAuthentication, 0)

	filter := bson.M{"status": models.StudentIdentityPhotoUploaded}
	cursor, err := db.GetCollection("auth").Find(context.TODO(), filter)
	if err != nil {
		return
	}

	if err = cursor.All(context.TODO(), &authProcs); err != nil {
		panic(err)
	}
	return
}

func (adminService *AdminService) UpdateStuIDAuthStatus(req *models.ModifyStuIDAuthStatRequest) (err error) {

	proc, err := adminService.GetStuIDAuthProc(req.AuthProcID)
	fmt.Println(proc)
	if err == nil {
		if req.Status == models.StudentIdentityAuthenticated {
			var hei *models.HEI
			var major *models.Major
			hei, err = GetHEIService().GetHEIByName(proc.SchoolName)
			if err != nil {
				return
			}
			major, err = GetMajorService().GetMajorByName(proc.MajorName)
			if err != nil {
				return
			}
			err = GetAuthService().UpdateUserSchoolInfo(req.UserID, &models.EduBGItem{
				HEIName:    hei.Name,
				HEIID:      hei.ID,
				MajorID:    major.ID,
				MajorName:  major.Name,
				Stage:      proc.Stage,
				AuthProcID: proc.ID,
			})
			if err != nil {
				return
			}
		}
		filter := bson.M{"userID": req.UserID}
		fmt.Println(req)
		err = db.GetCollection("auth").FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": bson.M{"status": req.Status, "suggestion": req.Suggestion}}).Err()
		if err != nil {
			return
		}
	}
	return
}

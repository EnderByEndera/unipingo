package services

import (
	"context"
	"melodie-site/server/db"
	"melodie-site/server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MajorService struct{}

var majorService *MajorService

func GetMajorService() *MajorService {
	if majorService == nil {
		majorService = &MajorService{}
	}
	return majorService
}

func (service *MajorService) GetMajorByName(majorName string) (major *models.Major, err error) {
	filter := bson.M{"name": majorName}
	err = db.GetCollection("majors").FindOne(context.TODO(), filter).Decode(&major)
	return
}

func (service *MajorService) GetMajor(majorID primitive.ObjectID) (major *models.Major, err error) {
	filter := bson.M{"_id": majorID}
	err = db.GetCollection("majors").FindOne(context.TODO(), filter).Decode(&major)
	return
}

// 添加一个通过专业大类来过滤专业的函数，参考上面的通过名称获取。

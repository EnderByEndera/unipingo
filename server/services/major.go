package services

import (
	"context"
	"melodie-site/server/db"
	"melodie-site/server/models"

	"go.mongodb.org/mongo-driver/bson"
)

type MajorService struct{}

var majorService *MajorService

func GetMajorService() *MajorService {
	if majorService == nil {
		majorService = &MajorService{}
	}
	return majorService
}

//通过名称查找major
func (service *MajorService) GetMajorByName(majorName string) (major *models.Major, err error) {
	filter := bson.M{"name": majorName}
	err = db.GetCollection("majors").FindOne(context.TODO(), filter).Decode(&major)
	return
}

//通过category进行筛选
//如果传入的值是"",那么返回所有
func (service *MajorService) FilterMajor(category string)(majors []*models.Major, err error){
	filter := bson.M{}
	if category!=""{
		filter["category"] = category
	}
	res, err := db.GetCollection("majors").Find(context.TODO(), filter)
	if err != nil {
		return
	}
	err = res.All(context.TODO(), &majors)
	return
}

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

func (service *MajorService) GetMajorByName(majorName string) (major *models.Major, err error) {
	filter := bson.M{"name": majorName}
	err = db.GetCollection("majors").FindOne(context.TODO(), filter).Decode(&major)
	return
}

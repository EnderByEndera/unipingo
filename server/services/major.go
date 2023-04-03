package services

import (
	"context"
	"melodie-site/server/db"
	"melodie-site/server/models"

	"go.mongodb.org/mongo-driver/bson"
)

type HEIService struct{}

var heiService *HEIService

func GetHEIService() *HEIService {
	if heiService == nil {
		heiService = &HEIService{}
	}
	return heiService
}

func (service *HEIService) GetHEIByName(heiName string) (hei *models.HEI, err error) {
	filter := bson.M{"name": heiName}
	err = db.GetCollection("heis").FindOne(context.TODO(), filter).Decode(&hei)
	return
}

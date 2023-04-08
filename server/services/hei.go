package services

import (
	"context"
	"melodie-site/server/db"
	"melodie-site/server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HEIService struct{}

var heiService *HEIService

func GetHEIService() *HEIService {
	if heiService == nil {
		heiService = &HEIService{}
	}
	return heiService
}

// 通过id获取学校
func (service *HEIService) GetHEI(id primitive.ObjectID) (hei *models.HEI, err error) {
	filter := bson.M{"_id": id}
	err = db.GetCollection("heis").FindOne(context.TODO(), filter).Decode(&hei)
	return
}

// 通过名称获取学校
// 如果error不为nil，说明没找到
func (service *HEIService) GetHEIByName(heiName string) (hei *models.HEI, err error) {
	filter := bson.M{"name": heiName}
	err = db.GetCollection("heis").FindOne(context.TODO(), filter).Decode(&hei)
	return
}

// 通过标签进行过滤
// 若要provincialLocation和tags不做限制，则传入"";
// 若要level和mode不做限制，则传入<0值
func (service *HEIService) FilterHEI(provincialLocation string, level models.HEILevel, mode models.HEIMode, policy string) (heis []*models.HEI, err error) {
	filter := bson.M{}
	if provincialLocation != "" {
		filter["location.provincial"] = provincialLocation
	}
	if level >= 0 {
		filter["level"] = level
	}
	if mode >= 0 {
		filter["mode"] = mode
	}
	if policy != "" {
		filter["tags"] = bson.M{
			"$elemMatch": bson.M{
				"$eq": policy,
			},
		}
	}
	res, err := db.GetCollection("heis").Find(context.TODO(), filter)
	if err != nil {
		return
	}
	err = res.All(context.TODO(), &heis)
	return
}

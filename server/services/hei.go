package services

import (
	"context"
	"melodie-site/server/db"
	"melodie-site/server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func (service *HEIService) FilterHEI(provincialLocation string, level models.HEILevel, mode models.HEIMode, policy string, page int64) (heis []*models.HEI, err error) {
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
	opts := options.Find().SetLimit(20).SetSkip(20 * page)
	res, err := db.GetCollection("heis").Find(context.TODO(), filter, opts)
	if err != nil {
		return
	}
	err = res.All(context.TODO(), &heis)
	return
}

func (service *HEIService) GetHEIName(majorID primitive.ObjectID) (name string, err error) {
	hei, err := service.GetHEI(majorID)
	if err != nil {
		return
	}
	name = hei.Name
	return
}

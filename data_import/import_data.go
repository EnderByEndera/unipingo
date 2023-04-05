package main

import (
	"context"
	"fmt"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"melodie-site/server/utils"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ImportMajors() {
	majors := make([]map[string]interface{}, 0)
	err := utils.UnmarshallJSONFromFile("majors.json", &majors)
	if err != nil {
		panic(err)
	}
	for i, majorProps := range majors {
		// majorProps

		major := models.Major{
			Code:     fmt.Sprintf("%+v", majorProps["专业代码"]),
			Name:     fmt.Sprintf("%+v", majorProps["专业名称"]),
			Category: fmt.Sprintf("%+v", majorProps["学位授予门类"]),
		}
		fmt.Printf("%d %+v\n", i, major)

		filter := bson.M{"name": major.Name}
		res := db.GetCollection("majors").FindOne(context.TODO(), filter)
		var insertedDocID primitive.ObjectID
		if res.Err() != nil {
			fmt.Println("err", res.Err())
			res, err := db.GetCollection("majors").InsertOne(context.TODO(), major)
			if err != nil {
				panic(err)
			}
			insertedDocID = res.InsertedID.(primitive.ObjectID)
		} else {
			_major := models.Major{}
			err = res.Decode(&_major)
			insertedDocID = _major.ID
			if err != nil {
				panic(err)
			}
		}
		fmt.Println(insertedDocID)

	}
}

func GetHEILevel(level string) models.HEILevel {
	if level == "本科" {
		return models.UniversityHEI
	} else {
		return models.VocationalHEI
	}
}

func GetHEIMode(mode string) models.HEIMode {
	if mode == "民办" {
		return models.PrivateHEI
	} else if mode == "内地与港澳台地区合作办学" {
		return models.Mainland_HK_MO_TW_Cooperative
	} else if mode == "中外合作办学" {
		return models.SinoForeignCooperative
	} else {
		return models.PublicHEI
	}
}

func ImportHEIs() {
	heis := make([]map[string]interface{}, 0)
	err := utils.UnmarshallJSONFromFile("heis.json", &heis)
	if err != nil {
		panic(err)
	}

	province := ""
	for i, majorProps := range heis {
		// majorProps
		if majorProps["code"] == nil {
			province = strings.Split(majorProps["序号"].(string), "(")[0]
			continue
		}
		hei := models.HEI{
			Code: fmt.Sprintf("%d", int(majorProps["code"].(float64))),
			Name: fmt.Sprintf("%+v", majorProps["name"]),
			Location: models.Location{
				Provincial: province,
				Municipal:  fmt.Sprintf("%+v", majorProps["region"]),
			},
			Level: GetHEILevel(fmt.Sprintf("%+v", majorProps["level"])),
			Mode:  GetHEIMode(fmt.Sprintf("%+v", majorProps["备注"])),
		}
		fmt.Printf("%d %+v\n", i, hei)

		filter := bson.M{"name": hei.Name}
		res := db.GetCollection("heis").FindOne(context.TODO(), filter)
		var insertedDocID primitive.ObjectID
		if res.Err() != nil {
			fmt.Println("err", res.Err())
			res, err := db.GetCollection("heis").InsertOne(context.TODO(), hei)
			if err != nil {
				panic(err)
			}
			insertedDocID = res.InsertedID.(primitive.ObjectID)
		} else {
			_major := models.Major{}
			err = res.Decode(&_major)
			insertedDocID = _major.ID
			if err != nil {
				panic(err)
			}
		}
		fmt.Println(insertedDocID)

	}
}

func main() {
	ImportMajors()
	ImportHEIs()
	fmt.Println("数据库初始化成功！")
}

package main

import (
	"context"
	"fmt"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"melodie-site/server/utils"

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

func main() {
	ImportMajors()
}

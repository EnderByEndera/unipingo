package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"melodie-site/server/db"
	"melodie-site/server/models"
)

var service *XZXJDiscussService

type XZXJDiscussService struct {
	wc         *writeconcern.WriteConcern
	rc         *readconcern.ReadConcern
	xzxjDBName string
	userDBName string
}

func GetXZXJDiscussService() *XZXJDiscussService {
	if service == nil {
		service = &XZXJDiscussService{
			wc:         writeconcern.New(writeconcern.W(1), writeconcern.J(true)),
			rc:         readconcern.New(readconcern.Level("local")),
			xzxjDBName: "xzxj_discuss",
			userDBName: "user",
		}
	}
	return service
}

func (service *XZXJDiscussService) AddOrUpdateXZXJUser(xzxjUserFormMeta *models.XZXJUserFormMeta) (xzxjID primitive.ObjectID, err error) {
	filter := bson.M{
		"userID": xzxjUserFormMeta.UserID,
	}

	opts := options.Update().SetUpsert(true)

	transaction := func(sessCtx mongo.SessionContext) (res interface{}, sessErr error) {
		_, sessErr = db.GetCollection(service.xzxjDBName).UpdateOne(sessCtx, filter, bson.M{"$setOnInsert": xzxjUserFormMeta.XZXJUser}, opts)
		if sessErr != nil {
			return
		}

		sessErr = db.GetCollection(service.userDBName).FindOneAndUpdate(sessCtx, bson.M{"_id": xzxjUserFormMeta.UserID}, bson.M{
			"$set": bson.M{
				"realName": xzxjUserFormMeta.RealName,
				"userTags": xzxjUserFormMeta.UserTags,
			}}).Err()
		if sessErr != nil {
			return
		}

		res = xzxjUserFormMeta.UserID

		return
	}
	sessOpts := options.Session().SetDefaultWriteConcern(service.wc).SetDefaultReadConcern(service.rc)
	result, err := db.GetMongoConn().UseSession(sessOpts, transaction)
	xzxjID, _ = result.(primitive.ObjectID)
	return
}

func (service *XZXJDiscussService) QueryXZXJUserByUserID(userID primitive.ObjectID) (xzxjUser *models.XZXJUser, err error) {
	filter := bson.M{"userID": userID}

	err = db.GetCollection(service.xzxjDBName).FindOne(context.TODO(), filter).Decode(xzxjUser)
	return
}

func (service *XZXJDiscussService) DeleteXZXJUserByID(userID primitive.ObjectID) (err error) {
	filter := bson.M{"userID": userID}

	err = db.GetCollection(service.xzxjDBName).FindOneAndDelete(context.TODO(), filter).Err()
	return
}

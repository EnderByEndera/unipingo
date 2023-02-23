package services

import (
	"context"
	"melodie-site/server/db"
	"melodie-site/server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var postService *PostService

type PostService struct {
	// PrivateKeys       map[uuid.UUID][]byte
	// WechatSessionKeys map[uuid.UUID]string
}

func (service *PostService) NewPost(post *models.Post) (docID primitive.ObjectID, err error) {
	conn := db.GetMongoConn()
	collection := conn.Client.Database("blog").Collection("posts")
	res1, err := collection.InsertOne(conn.Context, post)
	if err != nil {
		return
	}
	docID = res1.InsertedID.(primitive.ObjectID)
	return
}

func (service *PostService) GetPostByID(oid primitive.ObjectID) (post *models.Post, err error) {
	post = &models.Post{}
	conn := db.GetMongoConn()
	collection := conn.Client.Database("blog").Collection("posts")
	filter := bson.D{{"_id", oid}}
	err = collection.FindOne(context.TODO(), filter).Decode(post)
	return
}

func GetPostsService() *PostService {
	if postService == nil {
		postService = &PostService{}
	}
	return postService
}

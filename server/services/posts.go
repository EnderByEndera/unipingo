package services

import (
	"context"
	"errors"
	"fmt"
	"melodie-site/server/db"
	"melodie-site/server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var postService *PostService

type PostService struct {
	// PrivateKeys       map[uuid.UUID][]byte
	// WechatSessionKeys map[uuid.UUID]string
}

func getCollection(collectionName string) (collection *mongo.Collection) {
	conn := db.GetMongoConn()
	collection = conn.Client.Database("blog").Collection(collectionName)
	return
}

func (service *PostService) NewPost(post *models.Post) (docID primitive.ObjectID, err error) {
	conn := db.GetMongoConn()
	res1, err := getCollection("posts").InsertOne(conn.Context, post)
	if err != nil {
		return
	}
	docID = res1.InsertedID.(primitive.ObjectID)
	return
}

func (service *PostService) GetPostByID(oid primitive.ObjectID) (post *models.Post, err error) {
	post = &models.Post{}

	filter := bson.D{{"_id", oid}}
	err = getCollection("posts").FindOne(context.TODO(), filter).Decode(post)
	return
}

func (service *PostService) NewComment(req *models.NewCommentRequest) (commentUUID string, err error) {
	comment := models.NewComment(req)
	statement := bson.M{"$push": bson.M{"comments": comment}}
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)
	res := getCollection("posts").FindOneAndUpdate(context.TODO(), bson.D{{"_id", req.PostOID}}, statement, opts)
	err = res.Err()
	commentUUID = comment.UUID
	return
}

func (service *PostService) NewReply(req *models.NewReplyRequest) (newReplyUUID string, err error) {
	reply := models.NewReply(req)
	identifier := []interface{}{bson.D{{"commentItem.uuid", bson.D{{"$eq", req.CommentUUID}}}}}
	update := bson.D{{"$push", bson.D{{"comments.$[commentItem].replies", reply}}}}

	newReplyUUID = reply.UUID

	opts := options.FindOneAndUpdate().
		SetArrayFilters(options.ArrayFilters{Filters: identifier}).
		SetReturnDocument(options.After)
	var updatedDoc models.Post
	err = getCollection("posts").FindOneAndUpdate(context.TODO(), bson.D{{"_id", req.PostOID}}, update, opts).Decode(&updatedDoc)

	return
}

func (service *PostService) CheckIfAlreadyLiked(req *models.LikePostRequest) (alreadyLiked bool, alreadyDisLiked bool) {
	filter := bson.D{{"_id", req.PostOID}, {"likes", bson.M{"$elemMatch": bson.M{"userID": bson.M{"$eq": req.UserID}, "position": bson.M{"$eq": true}}}}}
	result := getCollection("posts").FindOne(context.TODO(), filter)
	if result.Err() == nil {
		alreadyLiked = true
	}

	filter = bson.D{{"_id", req.PostOID}, {"likes", bson.M{"$elemMatch": bson.M{"userID": bson.M{"$eq": req.UserID}, "position": bson.M{"$eq": false}}}}}
	result = getCollection("posts").FindOne(context.TODO(), filter)

	if result.Err() == nil {
		alreadyDisLiked = true
	}
	return
}

// 如果赞过了，就返回
func (service *PostService) GiveLikeToPost(req *models.LikePostRequest) (err error) {
	like := models.Like{UserID: req.UserID, TimeStamp: 0, Position: req.Position}
	alreadyLiked, alreadyDisliked := service.CheckIfAlreadyLiked(req)
	var statement bson.M
	identifier := []interface{}{}
	if !(alreadyLiked || alreadyDisliked) {
		if req.Position {
			statement = bson.M{"$push": bson.M{"likes": like}, "$inc": bson.M{"statistics.likes": 1}}
		} else {
			statement = bson.M{"$push": bson.M{"likes": like}, "$inc": bson.M{"statistics.dislikes": 1}}
		}
	} else if alreadyLiked {
		if req.Position {
			return errors.New("already liked this post")
		} else {
			// bson.M{"likeItem.userID": req.UserID},
			identifier = []interface{}{bson.D{{"likeItem.userID", bson.D{{"$eq", req.UserID}}}}}
			statement = bson.M{"$set": bson.M{"likes.$[likeItem].position": false},
				"$inc": bson.M{"statistics.likes": -1, "statistics.dislikes": 1}}
		}
	} else if alreadyDisliked {
		if !req.Position {
			return errors.New("already unliked this post")
		} else {
			identifier = []interface{}{bson.D{{"likeItem.userID", bson.D{{"$eq", req.UserID}}}}}
			statement = bson.M{"$set": bson.M{"likes.$[likeItem].position": true},
				"$inc": bson.M{"statistics.likes": 1, "statistics.dislikes": -1}}
		}
	} else {
		panic("error occurred. is liked and unliked exist at the same time?")
	}
	opts := options.FindOneAndUpdate().
		SetArrayFilters(options.ArrayFilters{Filters: identifier}).
		SetReturnDocument(options.After)

	res := getCollection("posts").FindOneAndUpdate(context.TODO(), bson.D{{"_id", req.PostOID}}, statement, opts)
	err = res.Err()
	return
}

func (service *PostService) CheckIfCommentAlreadyLiked(req *models.LikeCommentRequest) (alreadyLiked bool, alreadyDisLiked bool) {
	filter := bson.D{{"_id", req.PostOID}, {"comments",
		bson.M{"$elemMatch": bson.M{"uuid": bson.M{"$eq": req.CommentUUID},
			"likes": bson.M{"$elemMatch": bson.M{"userID": bson.M{"$eq": req.UserID}, "position": bson.M{"$eq": true}}}}}}}
	result := getCollection("posts").FindOne(context.TODO(), filter)
	if result.Err() == nil {
		alreadyLiked = true
	}

	filter = bson.D{{"_id", req.PostOID}, {"comments",
		bson.M{"$elemMatch": bson.M{"uuid": bson.M{"$eq": req.CommentUUID},
			"likes": bson.M{"$elemMatch": bson.M{"userID": bson.M{"$eq": req.UserID}, "position": bson.M{"$eq": false}}}}}}}
	result = getCollection("posts").FindOne(context.TODO(), filter)

	if result.Err() == nil {
		alreadyDisLiked = true
	}
	return
}

func (service *PostService) GiveLikeToComment(req *models.LikeCommentRequest) (err error) {
	like := models.Like{UserID: req.UserID, TimeStamp: 0, Position: req.Position}
	alreadyLiked, alreadyDisliked := service.CheckIfCommentAlreadyLiked(req)
	fmt.Println(alreadyLiked, alreadyDisliked)
	var statement bson.M
	identifier := []interface{}{}
	if !(alreadyLiked || alreadyDisliked) {
		if req.Position {
			identifier = []interface{}{bson.D{{"commentItem.uuid", bson.D{{"$eq", req.CommentUUID}}}}}
			statement = bson.M{"$push": bson.M{"comments.$[commentItem].likes": like}, "$inc": bson.M{"comments.$[commentItem].statistics.likes": 1}}
		} else {
			statement = bson.M{"$push": bson.M{"likes": like}, "$inc": bson.M{"statistics.dislikes": 1}}
		}
	} else if alreadyLiked {
		if req.Position {
			return errors.New("already liked this post")
		} else {
			// bson.M{"likeItem.userID": req.UserID},
			identifier = []interface{}{bson.D{{"likeItem.userID", bson.D{{"$eq", req.UserID}}}}}
			statement = bson.M{"$set": bson.M{"likes.$[likeItem].position": false},
				"$inc": bson.M{"statistics.likes": -1, "statistics.dislikes": 1}}
		}
	} else if alreadyDisliked {
		if !req.Position {
			return errors.New("already unliked this post")
		} else {
			identifier = []interface{}{bson.D{{"likeItem.userID", bson.D{{"$eq", req.UserID}}}}}
			statement = bson.M{"$set": bson.M{"likes.$[likeItem].position": true},
				"$inc": bson.M{"statistics.likes": 1, "statistics.dislikes": -1}}
		}
	} else {
		panic("error occurred. is liked and unliked exist at the same time?")
	}
	opts := options.FindOneAndUpdate().
		SetArrayFilters(options.ArrayFilters{Filters: identifier}).
		SetReturnDocument(options.After)

	res := getCollection("posts").FindOneAndUpdate(context.TODO(), bson.D{{"_id", req.PostOID}}, statement, opts)
	err = res.Err()
	return
}

func GetPostsService() *PostService {
	if postService == nil {
		postService = &PostService{}
	}
	return postService
}

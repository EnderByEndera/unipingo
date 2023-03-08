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

func (service *PostService) GetAllUserPosts(userid primitive.ObjectID) (posts []models.Post, err error) {
	posts = make([]models.Post, 0)

	filter := bson.D{{"userID", userid}}
	cursor, err := getCollection("posts").Find(context.TODO(), filter)
	if err != nil {
		return
	}

	if err = cursor.All(context.TODO(), &posts); err != nil {
		panic(err)
	}

	return
}

func (service *PostService) NewComment(req *models.NewCommentRequest) (commentUUID string, err error) {
	comment, err := models.NewComment(req)
	if err != nil {
		return
	}
	fmt.Println(req.PostOID)
	statement := bson.M{"$push": bson.M{"comments": comment}}
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)
	postOID, err := primitive.ObjectIDFromHex(req.PostOID)
	if err != nil {
		return
	}
	res := getCollection("posts").FindOneAndUpdate(context.TODO(), bson.M{"_id": postOID}, statement, opts)
	err = res.Err()
	commentUUID = comment.UUID
	return
}

func (service *PostService) NewReply(req *models.NewReplyRequest) (newReplyUUID string, err error) {
	reply := models.NewReply(req)
	postOID, err := primitive.ObjectIDFromHex(req.PostOID)
	if err != nil {
		return
	}
	identifier := []interface{}{bson.D{{"commentItem.uuid", bson.D{{"$eq", req.CommentUUID}}}}}
	update := bson.D{{"$push", bson.D{{"comments.$[commentItem].replies", reply}}}}

	newReplyUUID = reply.UUID

	opts := options.FindOneAndUpdate().
		SetArrayFilters(options.ArrayFilters{Filters: identifier}).
		SetReturnDocument(options.After)
	var updatedDoc models.Post
	err = getCollection("posts").FindOneAndUpdate(context.TODO(), bson.M{"_id": postOID}, update, opts).Decode(&updatedDoc)

	return
}

func exists(filter bson.M) bool {
	result := getCollection("posts").FindOne(context.TODO(), filter)
	return result.Err() == nil
}

func (service *PostService) CheckIfAlreadyLiked(postOID, userID primitive.ObjectID) (alreadyLiked bool, alreadyDisLiked bool) {
	filterCreator := func(position bool) bson.M {

		return bson.M{
			"_id": postOID,
			"likes": bson.M{
				"$elemMatch": bson.M{
					"userID": bson.M{
						"$eq": userID,
					},
					"position": bson.M{
						"$eq": position,
					},
				},
			},
		}
	}
	return exists(filterCreator(true)), exists(filterCreator(false))
}

// 如果赞过了，就返回
func (service *PostService) GiveLikeToPost(req *models.LikePostRequest) (err error) {
	postOID, err := primitive.ObjectIDFromHex(req.PostOID)
	if err != nil {
		return
	}
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return
	}
	like := models.Like{UserID: userID, TimeStamp: 0, Position: req.Position}
	alreadyLiked, alreadyDisliked := service.CheckIfAlreadyLiked(postOID, userID)
	var statement bson.M
	identifier := []interface{}{}
	incFactory := func(likeInc, dislikeInc int) bson.M {
		return bson.M{"statistics.likes": likeInc, "statistics.dislikes": dislikeInc}
	}
	modificationFactory := func(position bool, likeInc, dislikeInc int) ([]interface{}, bson.M) {
		return []interface{}{bson.M{
				"likeItem.userID": bson.M{
					"$eq": userID,
				},
			}}, bson.M{
				"$set": bson.M{
					"likes.$[likeItem].position": position,
				},
				"$inc": incFactory(likeInc, dislikeInc)}
	}
	if !(alreadyLiked || alreadyDisliked) {
		if req.Position {
			statement = bson.M{"$push": bson.M{"likes": like}, "$inc": incFactory(1, 0)}
		} else {
			statement = bson.M{"$push": bson.M{"likes": like}, "$inc": incFactory(0, 1)}
		}
	} else if alreadyLiked {
		if req.Position {
			return errors.New("already liked this post")
		} else {
			identifier, statement = modificationFactory(false, -1, 1)
		}
	} else if alreadyDisliked {
		if !req.Position {
			return errors.New("already unliked this post")
		} else {
			identifier, statement = modificationFactory(true, 1, -1)
		}
	} else {
		panic("error occurred. is liked and unliked exist at the same time?")
	}
	opts := options.FindOneAndUpdate().
		SetArrayFilters(options.ArrayFilters{Filters: identifier}).
		SetReturnDocument(options.After)

	res := getCollection("posts").FindOneAndUpdate(context.TODO(), bson.D{{"_id", postOID}}, statement, opts)
	err = res.Err()
	return
}

func (service *PostService) CheckIfCommentAlreadyLiked(postOID primitive.ObjectID, userID primitive.ObjectID, commentUUID string) (alreadyLiked bool, alreadyDisLiked bool) {
	filterCreator := func(pos bool) bson.M {
		return bson.M{
			"_id": postOID,
			"comments": bson.M{
				"$elemMatch": bson.M{
					"uuid": bson.M{
						"$eq": commentUUID,
					},
					"likes": bson.M{
						"$elemMatch": bson.M{
							"userID": bson.M{
								"$eq": userID,
							},
							"position": bson.M{"$eq": pos},
						},
					},
				},
			},
		}
	}

	return exists(filterCreator(true)), exists(filterCreator(false))
}

func (service *PostService) GiveLikeToComment(req *models.LikeCommentRequest) (err error) {
	postOID, err := primitive.ObjectIDFromHex(req.PostOID)
	if err != nil {
		return
	}
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return
	}
	like := models.Like{UserID: userID, TimeStamp: 0, Position: req.Position}
	alreadyLiked, alreadyDisliked := service.CheckIfCommentAlreadyLiked(postOID, userID, req.CommentUUID)

	var statement bson.M
	identifier := []interface{}{}
	pushFactory := func(likeInc, dislikeInc int) ([]interface{}, bson.M) {
		identifier_ := []interface{}{
			bson.M{
				"commentItem.uuid": bson.M{
					"$eq": req.CommentUUID,
				},
			},
		}
		statement_ := bson.M{
			"$push": bson.M{"comments.$[commentItem].likes": like},
			"$inc": bson.M{
				"comments.$[commentItem].statistics.likes":    likeInc,
				"comments.$[commentItem].statistics.dislikes": dislikeInc,
			},
		}
		return identifier_, statement_
	}

	modificationFactory := func(position bool, likeInc, dislikeInc int) ([]interface{}, bson.M) {
		identifier_ := []interface{}{
			bson.M{"commentItem.uuid": bson.M{
				"$eq": req.CommentUUID,
			}},
			bson.M{"likeItem.userID": bson.M{
				"$eq": userID,
			}},
		}
		statement_ := bson.M{
			"$set": bson.M{
				"comments.$[commentItem].likes.$[likeItem].position": position,
			},
			"$inc": bson.M{
				"comments.$[commentItem].statistics.dislikes": likeInc,
				"comments.$[commentItem].statistics.likes":    dislikeInc,
			},
		}
		return identifier_, statement_
	}
	if !(alreadyLiked || alreadyDisliked) {
		if req.Position {
			identifier, statement = pushFactory(1, 0)
		} else {
			identifier, statement = pushFactory(0, 1)
		}
	} else if alreadyLiked {
		if req.Position {
			return errors.New("already liked this post")
		} else {
			identifier, statement = modificationFactory(false, -1, 1)
		}
	} else if alreadyDisliked {
		if !req.Position {
			return errors.New("already disliked this post")
		} else {
			identifier, statement = modificationFactory(true, 1, -1)
		}
	} else {
		panic("error occurred. is liked and unliked exist at the same time?")
	}
	opts := options.FindOneAndUpdate().
		SetArrayFilters(options.ArrayFilters{Filters: identifier}).
		SetReturnDocument(options.After)

	res := getCollection("posts").FindOneAndUpdate(context.TODO(), bson.D{{"_id", postOID}}, statement, opts)
	err = res.Err()
	return
}

func (service *PostService) CheckIfReplyAlreadyLiked(postOID primitive.ObjectID, userID primitive.ObjectID, commentUUID, replyUUID string) (alreadyLiked bool, alreadyDisLiked bool) {

	createFilter := func(position bool) bson.M {
		v := bson.M{
			"_id": postOID,
			"comments": bson.M{
				"$elemMatch": bson.M{
					"uuid": bson.M{"$eq": commentUUID},
					"replies": bson.M{
						"$elemMatch": bson.M{
							"uuid": bson.M{"$eq": replyUUID},
							"likes": bson.M{
								"$elemMatch": bson.M{
									"userID":   bson.M{"$eq": userID},
									"position": bson.M{"$eq": position},
								},
							},
						},
					},
				},
			},
		}
		return v
	}
	return exists(createFilter(true)), exists(createFilter(false))
}

func (service *PostService) GiveLikeToReply(req *models.LikeReplyRequest) (err error) {
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return
	}
	postID, err := primitive.ObjectIDFromHex(req.PostOID)
	if err != nil {
		return
	}
	like := models.Like{UserID: userID, TimeStamp: 0, Position: req.Position}
	alreadyLiked, alreadyDisliked := service.CheckIfReplyAlreadyLiked(postID, userID, req.CommentUUID, req.ReplyUUID)
	var statement bson.M
	identifier := []interface{}{}
	pushFactory := func(likeInc, dislikeInc int) ([]interface{}, bson.M) {
		identifier_ := []interface{}{
			bson.M{
				"commentItem.uuid": bson.M{
					"$eq": req.CommentUUID,
				},
			},
			bson.M{
				"replyItem.uuid": bson.M{
					"$eq": req.ReplyUUID,
				},
			},
		}
		statement_ := bson.M{
			"$push": bson.M{
				"comments.$[commentItem].replies.$[replyItem].likes": like,
			},
			"$inc": bson.M{
				"comments.$[commentItem].replies.$[replyItem].statistics.likes":    likeInc,
				"comments.$[commentItem].replies.$[replyItem].statistics.dislikes": dislikeInc,
			},
		}
		return identifier_, statement_
	}

	modificationFactory := func(position bool, likeInc, dislikeInc int) ([]interface{}, bson.M) {
		identifier_ := []interface{}{
			bson.M{
				"commentItem.uuid": bson.M{
					"$eq": req.CommentUUID,
				},
			},
			bson.M{
				"replyItem.uuid": bson.M{
					"$eq": req.ReplyUUID,
				},
			},
			bson.M{
				"likeItem.userID": bson.M{
					"$eq": userID,
				},
			},
		}
		statement_ := bson.M{
			"$set": bson.M{
				"comments.$[commentItem].replies.$[replyItem].likes.$[likeItem].position": false,
			},
			"$inc": bson.M{
				"comments.$[commentItem].replies.$[replyItem].statistics.dislikes": 1,
				"comments.$[commentItem].replies.$[replyItem].statistics.likes":    -1,
			},
		}
		return identifier_, statement_
	}
	if !(alreadyLiked || alreadyDisliked) {
		if req.Position {
			identifier, statement = pushFactory(1, 0)
		} else {
			identifier, statement = pushFactory(0, 1)
		}
	} else if alreadyLiked {
		if req.Position {
			return errors.New("already liked this post")
		} else {
			identifier, statement = modificationFactory(false, -1, 1)
		}
	} else if alreadyDisliked {
		if !req.Position {
			return errors.New("already disliked this post")
		} else {
			identifier, statement = modificationFactory(true, 1, -1)
		}
	} else {
		panic("error occurred. is liked and unliked exist at the same time?")
	}
	opts := options.FindOneAndUpdate().
		SetArrayFilters(options.ArrayFilters{Filters: identifier}).
		SetReturnDocument(options.After)

	res := getCollection("posts").FindOneAndUpdate(context.TODO(), bson.D{{"_id", postID}}, statement, opts)
	err = res.Err()
	return
}

func (service *PostService) PostToOutline(post *models.Post) (outline models.PostOutline, err error) {
	user, err := authService.GetUserByID(post.UserID)
	if err != nil {
		return
	}
	outline = models.PostOutline{
		PostOID:    post.DocumentID.Hex(),
		TimeStamp:  post.TimeStamp,
		User:       user.ToPublicInfo(),
		Title:      post.Title,
		Content:    post.Content,
		Statistics: post.Statistics,
	}
	return outline, nil
}

func GetPostsService() *PostService {
	if postService == nil {
		postService = &PostService{}
	}
	return postService
}

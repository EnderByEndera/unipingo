package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"testing"

	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func TestModelPosts(t *testing.T) {
	byts, _ := ioutil.ReadFile("post.json")
	post := models.Post{}
	err := json.Unmarshal(byts, &post)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	res, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Println(string(res))

	mongoConn := db.GetMongoConn()
	client := mongoConn.Client

	collection := client.Database("blog").Collection("posts")
	insertedDocID, err := services.GetPostsService().NewPost(&post)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	postGot, err := services.GetPostsService().GetPostByID(insertedDocID)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if postGot.DocumentID != insertedDocID {
		fmt.Println("Document id got not equal！")
		t.FailNow()
	}
	//////////////////////////////////
	// filter := bson.D{{"_id", insertedDocID}}
	// update := bson.D{{"$set", bson.D{{"title", "Ledebouria socialis"}}}}
	// _, err = collection.UpdateOne(context.TODO(), filter, update)
	// if err != nil {
	// 	panic(err)
	// }
	// u, _ := uuid.Parse("de867ca7-270e-4b00-a7d6-37bd8f073058")
	// identifier := []interface{}{bson.D{{"hotOptions.uuid", bson.D{{"$eq", u}}}}}
	// update = bson.D{{"$set", bson.D{{"comments.$[hotOptions].content", "哈哈哈哈哈哈hahaqqqhda"}}}}
	// opts := options.FindOneAndUpdate().
	// 	SetArrayFilters(options.ArrayFilters{Filters: identifier}).
	// 	SetReturnDocument(options.After)
	// var updatedDoc models.Post
	// err = collection.FindOneAndUpdate(context.TODO(), bson.D{{"_id", insertedDocID}}, update, opts).Decode(&updatedDoc)
	//////////////////////////////////////////////////////
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("=================\n", updatedDoc.Comments)

	// res, _ = bson.MarshalExtJSON(updatedDoc, false, false)
	// fmt.Println(string(res))

	// statement := bson.M{"$pull": bson.M{"userActions": bson.M{"type": bson.D{{"$eq", 1}}, "userID": bson.D{{"$eq", 2}}}}}
	// opts = options.FindOneAndUpdate().
	// 	SetReturnDocument(options.After)
	// err = collection.FindOneAndUpdate(context.TODO(), bson.D{{"_id", insertedDocID}}, statement, opts).Decode(&updatedDoc)

	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("=================\n", updatedDoc.UserActions)

	// res, _ = bson.MarshalExtJSON(updatedDoc, false, false)
	// fmt.Println(string(res))

	commentUUID, err := services.GetPostsService().NewComment(&models.NewCommentRequest{PostOID: insertedDocID, UserID: 77, Content: "日内瓦，退钱！"})
	if err != nil {
		panic(err)
	}
	replyID, err := services.GetPostsService().NewReply(&models.NewReplyRequest{PostOID: insertedDocID, CommentUUID: commentUUID, UserID: 999, Content: "赞同，暴躁老哥！", ToUUID: commentUUID})
	if err != nil {
		panic(err)
	}
	fmt.Println(replyID)

	err = services.GetPostsService().GiveLikeToPost(&models.LikePostRequest{PostOID: insertedDocID, UserID: 88, Position: true})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	postGot, err = services.GetPostsService().GetPostByID(insertedDocID)
	assert.Equal(t, postGot.Likes[0].UserID, 88)
	assert.Equal(t, postGot.Likes[0].Position, true)
	assert.Equal(t, postGot.Likes[0].Position, true)
	assert.Equal(t, postGot.Statistics.Likes, 1)
	assert.Equal(t, postGot.Statistics.Dislikes, 0)

	err = services.GetPostsService().GiveLikeToPost(&models.LikePostRequest{PostOID: insertedDocID, UserID: 88, Position: false})

	postGot, err = services.GetPostsService().GetPostByID(insertedDocID)
	assert.Equal(t, postGot.Likes[0].UserID, 88)
	assert.Equal(t, postGot.Likes[0].Position, false)
	assert.Equal(t, postGot.Statistics.Likes, 0)
	assert.Equal(t, postGot.Statistics.Dislikes, 1)

	err = services.GetPostsService().GiveLikeToPost(&models.LikePostRequest{PostOID: insertedDocID, UserID: 88, Position: true})

	postGot, err = services.GetPostsService().GetPostByID(insertedDocID)
	assert.Equal(t, postGot.Likes[0].UserID, 88)
	assert.Equal(t, postGot.Likes[0].Position, true)
	assert.Equal(t, postGot.Statistics.Likes, 1)
	assert.Equal(t, postGot.Statistics.Dislikes, 0)

	err = services.GetPostsService().GiveLikeToComment(&models.LikeCommentRequest{
		PostOID: insertedDocID, UserID: 88, Position: true, CommentUUID: "de867ca7-270e-4b00-a7d6-37bd8f073058",
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	postGot, err = services.GetPostsService().GetPostByID(insertedDocID)
	assert.Equal(t, postGot.Comments[0].Likes[0].UserID, 88)
	assert.Equal(t, postGot.Comments[0].Likes[0].Position, true)
	err = services.GetPostsService().GiveLikeToComment(&models.LikeCommentRequest{
		PostOID: insertedDocID, UserID: 88, Position: false, CommentUUID: "de867ca7-270e-4b00-a7d6-37bd8f073058",
	})
	fmt.Println(err)

	postGot, err = services.GetPostsService().GetPostByID(insertedDocID)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	assert.Equal(t, postGot.Comments[0].Likes[0].UserID, 88)
	assert.Equal(t, postGot.Comments[0].Likes[0].Position, false)
	fmt.Println(postGot.ToIndentedJSON())
	res11, err := collection.DeleteOne(context.TODO(), bson.D{{"_id", insertedDocID}})
	if err != nil {
		panic(err)
	}
	fmt.Println(res11)

}

package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var UserIDHex string

func GiveLikeToPost(t *testing.T, insertedDocID primitive.ObjectID) (err error) {
	err = services.GetPostsService().GiveLikeToPost(&models.LikePostRequest{PostOID: insertedDocID.Hex(), UserID: UserIDHex, Position: true})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	postGot, err := services.GetPostsService().GetPostByID(insertedDocID)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	assert.Equal(t, postGot.Likes[0].UserID.Hex(), UserIDHex)
	assert.Equal(t, postGot.Likes[0].Position, true)
	assert.Equal(t, postGot.Likes[0].Position, true)
	assert.Equal(t, postGot.Statistics.Likes, 1)
	assert.Equal(t, postGot.Statistics.Dislikes, 0)

	err = services.GetPostsService().GiveLikeToPost(&models.LikePostRequest{PostOID: insertedDocID.Hex(), UserID: UserIDHex, Position: false})
	if err != nil {
		return
	}
	postGot, err = services.GetPostsService().GetPostByID(insertedDocID)
	if err != nil {
		return
	}
	assert.Equal(t, postGot.Likes[0].UserID.Hex(), UserIDHex)
	assert.Equal(t, postGot.Likes[0].Position, false)
	assert.Equal(t, postGot.Statistics.Likes, 0)
	assert.Equal(t, postGot.Statistics.Dislikes, 1)

	err = services.GetPostsService().GiveLikeToPost(&models.LikePostRequest{PostOID: insertedDocID.Hex(), UserID: UserIDHex, Position: true})
	if err != nil {
		return
	}
	postGot, err = services.GetPostsService().GetPostByID(insertedDocID)
	if err != nil {
		return
	}
	assert.Equal(t, postGot.Likes[0].UserID.Hex(), UserIDHex)
	assert.Equal(t, postGot.Likes[0].Position, true)
	assert.Equal(t, postGot.Statistics.Likes, 1)
	assert.Equal(t, postGot.Statistics.Dislikes, 0)
	return
}

func GiveLikeToComment(t *testing.T, insertedDocID primitive.ObjectID) (err error) {
	err = services.GetPostsService().GiveLikeToComment(&models.LikeCommentRequest{
		PostOID: insertedDocID.Hex(), UserID: UserIDHex, Position: true, CommentUUID: "de867ca7-270e-4b00-a7d6-37bd8f073058",
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	postGot, err := services.GetPostsService().GetPostByID(insertedDocID)
	assert.Equal(t, postGot.Comments[0].Likes[0].UserID.Hex(), UserIDHex)
	assert.Equal(t, postGot.Comments[0].Likes[0].Position, true)
	assert.Equal(t, postGot.Statistics.Likes, 1)
	assert.Equal(t, postGot.Statistics.Dislikes, 0)

	err = services.GetPostsService().GiveLikeToComment(&models.LikeCommentRequest{
		PostOID: insertedDocID.Hex(), UserID: UserIDHex, Position: false, CommentUUID: "de867ca7-270e-4b00-a7d6-37bd8f073058",
	})
	fmt.Println(err)

	postGot, err = services.GetPostsService().GetPostByID(insertedDocID)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	assert.Equal(t, postGot.Comments[0].Likes[0].UserID.Hex(), UserIDHex)
	assert.Equal(t, postGot.Comments[0].Likes[0].Position, false)
	fmt.Println(postGot.ToIndentedJSON())
	return
}

func GiveLikeToReply(t *testing.T, insertedDocID primitive.ObjectID) (err error) {

	// a, b := services.GetPostsService().CheckIfReplyAlreadyLiked()
	err = services.GetPostsService().GiveLikeToReply(&models.LikeReplyRequest{
		PostOID:     insertedDocID.Hex(),
		CommentUUID: "de867ca7-270e-4b00-a7d6-37bd8f073058",
		ReplyUUID:   "24c73dfb-8666-4d92-af84-af81a8211e97",
		UserID:      UserIDHex,
		Position:    true,
	})
	if err != nil {
		panic(err)
	}
	postGot, err := services.GetPostsService().GetPostByID(insertedDocID)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	assert.Equal(t, postGot.Comments[0].Replies[0].Likes[0].Position, true)
	assert.Equal(t, postGot.Comments[0].Replies[0].Likes[0].UserID.Hex(), UserIDHex)

	err = services.GetPostsService().GiveLikeToReply(&models.LikeReplyRequest{
		PostOID:     insertedDocID.Hex(),
		CommentUUID: "de867ca7-270e-4b00-a7d6-37bd8f073058",
		ReplyUUID:   "24c73dfb-8666-4d92-af84-af81a8211e97",
		UserID:      UserIDHex,
		Position:    false,
	})
	if err != nil {
		panic(err)
	}
	postGot, err = services.GetPostsService().GetPostByID(insertedDocID)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	assert.Equal(t, postGot.Comments[0].Replies[0].Likes[0].UserID.Hex(), UserIDHex)
	assert.Equal(t, postGot.Comments[0].Replies[0].Likes[0].Position, false)
	assert.Equal(t, postGot.Comments[0].Replies[0].Statistics.Likes, 0)
	assert.Equal(t, postGot.Comments[0].Replies[0].Statistics.Dislikes, 1)
	return
}

func TestModelPosts(t *testing.T) {
	user, err := services.GetAuthService().GetUserByName("admin")
	if err != nil {
		panic(err)
	}
	UserIDHex = user.OID.Hex()
	byts, _ := ioutil.ReadFile("post.json")
	post := models.Post{UserID: user.OID, TimeStamp: uint64(time.Now().UnixMilli())}
	err = json.Unmarshal(byts, &post)
	if err != nil {
		fmt.Println("marshal failed:", err)
		t.FailNow()
	}
	// res, err := json.MarshalIndent(post, "", "  ")
	// if err != nil {
	// 	fmt.Println(err)
	// 	t.FailNow()
	// }

	// fmt.Println(string(res))

	mongoConn := db.GetMongoConn()
	client := mongoConn.Client

	collection := client.Database("blog").Collection("posts")
	insertedDocID, err := services.GetPostsService().NewPost(&post)
	if err != nil {
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

	// assert.Equal(insertedDocID.Hex()==)
	commentUUID, err := services.GetPostsService().NewComment(&models.NewCommentRequest{PostOID: insertedDocID.Hex(), UserID: UserIDHex, Content: "日内瓦，退钱！"})
	assert.Equal(t, err, nil)

	replyID, err := services.GetPostsService().NewReply(&models.NewReplyRequest{PostOID: insertedDocID.Hex(), CommentUUID: commentUUID, UserID: UserIDHex, Content: "赞同，暴躁老哥！", ToUUID: commentUUID})
	if err != nil {
		panic(err)
	}
	fmt.Println(replyID)

	err = GiveLikeToPost(t, insertedDocID)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	err = GiveLikeToComment(t, insertedDocID)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	err = GiveLikeToReply(t, insertedDocID)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	// res11, err := collection.DeleteOne(context.TODO(), bson.D{{"_id", insertedDocID}})
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Println(collection, UserIDHex)

}

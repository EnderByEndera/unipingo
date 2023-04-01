package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

var userIDHex string

func TestModelPosts(t *testing.T) {
	user, err := services.GetAuthService().GetUserByName("admin")
	if err != nil {
		panic(err)
	}
	userIDHex = user.OID.Hex()
	byts, _ := ioutil.ReadFile("post.json")
	post := models.Post{UserID: user.OID, CreateTime: uint64(time.Now().UnixMilli())}
	err = json.Unmarshal(byts, &post)
	if err != nil {
		fmt.Println("marshal failed:", err)
		t.FailNow()
	}
	insertedDocID, err := services.GetPostsService().NewPost(&post)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(insertedDocID)

	postPtr, err := services.GetPostsService().GetPostByID(insertedDocID)
	assert.Equal(t, err, nil)
	fmt.Printf("%+v\n", postPtr)
}

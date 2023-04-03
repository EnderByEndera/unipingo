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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userIDHex string

var userName = "user-1"

var tempHEIID, _ = primitive.ObjectIDFromHex("642996289972d48e335fb0a5")

func TestModelPosts(t *testing.T) {
	user, err := services.GetAuthService().InternalAddUser(userName, "123456", models.RoleUnpaidUser, func(u *models.User) {
		u.EducationalBackground = make([]models.EduBGItem, 0)
		u.EducationalBackground = append(u.EducationalBackground, models.EduBGItem{HEIID: tempHEIID})
	})
	assert.Equal(t, err, nil)
	isAlumn, err := services.GetAuthService().IsAlumn(user.ID, tempHEIID)

	fmt.Printf("%+v\n", user)
	assert.Equal(t, isAlumn, true)
	assert.Equal(t, err, nil)
	// user, err := services.GetAuthService().GetUserByName("admin")
	if err != nil {
		panic(err)
	}
	userIDHex = user.ID.Hex()
	byts, _ := ioutil.ReadFile("post.json")
	major, err := services.GetMajorService().GetMajorByName("生物医学工程")
	assert.Equal(t, err, nil)
	post := models.Answer{
		UserID:     user.ID,
		CreateTime: uint64(time.Now().UnixMilli()),
		// Category: ,
		BelongsTo: models.EntityWithName{Name: major.Name, ID: major.ID},
	}

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

	err = services.GetPostsService().GiveLikeToPost(user.ID, insertedDocID)
	assert.Equal(t, err, nil)

	err = services.GetPostsService().GiveLikeToPost(user.ID, insertedDocID)

	postPtr, err := services.GetPostsService().GetPostByID(insertedDocID)
	assert.Equal(t, err, nil)
	fmt.Printf("%+v\n", postPtr)

	services.GetPostsService().DeletePostByID(insertedDocID)

	err = services.GetAuthService().InternalRemoveUser(userName)
	assert.Equal(t, err, nil)

	_, err = services.GetPostsService().GetPostByID(insertedDocID)
	assert.NotEqual(t, err, nil)
}

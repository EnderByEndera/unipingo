package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"
)

func testPassword(t *testing.T) {
	_, code, _ := services.GetAuthService().Login("user1", "123456")
	assert.Equal(t, code, http.StatusOK)
	_, code, _ = services.GetAuthService().Login("user1", "123457")
	assert.Equal(t, code, http.StatusBadRequest)
}

func testChangePublicInfo(t *testing.T) {
	user, err := services.GetAuthService().GetUserByName("user1")
	assert.Equal(t, err, nil)
	// newName := "超超威蓝猫"
	newAvatar := "1145.jpg"
	err = services.GetAuthService().UpdateUserPublicInfo(user.ID, &models.UserPublicInfoUpdateRequest{
		Avatar: newAvatar,
		Name:   "user1",
	})
	assert.Equal(t, err, nil)
	user, err = services.GetAuthService().GetUserByName("user1")
	assert.Equal(t, nil, err)
	assert.Equal(t, user.Avatar, newAvatar)
}

func testCollection(t *testing.T) {
	user, err := services.GetAuthService().GetUserByName("user1")
	assert.Equal(t, err, nil)

	hei, err := services.GetHEIService().GetHEIByName("北京航空航天大学")
	assert.Equal(t, err, nil)

	major, err := services.GetMajorService().GetMajorByName("能源与动力工程")
	assert.Equal(t, err, nil)
	fmt.Println(major)

	ret, err := services.GetAuthService().IsHEIOrMajorInCollection(user.ID, hei.ID, models.CollectionItemHEI)
	assert.Equal(t, ret, false)

	err = services.GetAuthService().AddHEIOrMajorToCollection(user.ID, hei.ID, models.CollectionItemHEI)
	assert.Equal(t, err, nil)

	ret, err = services.GetAuthService().IsHEIOrMajorInCollection(user.ID, hei.ID, models.CollectionItemHEI)
	assert.Equal(t, ret, true)

	err = services.GetAuthService().RemoveHEIOrMajorFromCollection(user.ID, hei.ID, models.CollectionItemHEI)
	assert.Equal(t, err, nil)

	ret, err = services.GetAuthService().IsHEIOrMajorInCollection(user.ID, hei.ID, models.CollectionItemHEI)
	assert.Equal(t, ret, false)

	err = services.GetAuthService().AddHEIOrMajorToCollection(user.ID, major.ID, models.CollectionItemMajor)
	assert.Equal(t, err, nil)

	user, err = services.GetAuthService().GetUserByName("user1")
	assert.Equal(t, err, nil)
	fmt.Printf("%+v\n", user.Collection)
	// assert.Equal()
}

func TestUser(t *testing.T) {
	user, err := services.GetAuthService().InternalAddUser("user1", "123456", models.RoleUnpaidUser, nil)
	assert.Equal(t, err, nil)
	byts, err := json.Marshal(user)
	assert.Equal(t, err, nil)
	user1 := models.User{}
	err = json.Unmarshal(byts, &user1)
	fmt.Println(err)
	assert.Equal(t, err, nil)
	fmt.Println("unmarshalled", user1)
	testPassword(t)
	testChangePublicInfo(t)
	testCollection(t)
	err = services.GetAuthService().InternalRemoveUser("user1")
	fmt.Println(err)
	_, err = services.GetAuthService().GetUserByName("user1")
	assert.NotEqual(t, err, nil)

}

func TestGetUserTag(t *testing.T) {
	admin, err := services.GetAuthService().GetUserByName("admin")
	fmt.Println(admin)
	assert.Equal(t, err, nil)
	tags, err := services.GetAuthService().GetTagByUserID(admin.ID)
	fmt.Print("tags是")
	fmt.Println(*tags)
	assert.Equal(t, err, nil)

	userID := primitive.NilObjectID
	_, err = services.GetAuthService().GetTagByUserID(userID)
	assert.Equal(t, err, errors.New("用户ID为空"))

	userID, err = primitive.ObjectIDFromHex("6464decb74edfbd3023d19d3")
	assert.Equal(t, err, nil)
	_, err = services.GetAuthService().GetTagByUserID(userID)
	assert.Equal(t, err, errors.New("数据库查找失败"))

}

func TestUpdateUserTag(t *testing.T) {
	admin, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)
	tags := []string{"喜欢理科"}
	err = services.GetAuthService().UpdateUserTag(admin.ID, tags)
	assert.Equal(t, err, nil)

	userID := primitive.NilObjectID
	//assert.Equal(t, err, nil)
	err = services.GetAuthService().UpdateUserTag(userID, tags)
	assert.Equal(t, err, errors.New("用户ID为空"))

}

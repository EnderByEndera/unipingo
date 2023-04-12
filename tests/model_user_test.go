package tests

import (
	"encoding/json"
	"fmt"
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

	err = services.GetAuthService().AddHEIOrMajorToCollection(user.ID, hei.ID, models.CollectionItemHEI)
	assert.Equal(t, err, nil)

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

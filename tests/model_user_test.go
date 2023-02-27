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
	_, code, _ := services.GetAuthService().Login("admin", "123456")
	assert.Equal(t, code, http.StatusOK)
	_, code, _ = services.GetAuthService().Login("admin", "123457")
	assert.Equal(t, code, http.StatusBadRequest)
}

func TestUser(t *testing.T) {
	user, err := services.GetAuthService().InternalAddUser("admin", "123456")
	assert.Equal(t, err, nil)
	byts, err := json.Marshal(user)
	assert.Equal(t, err, nil)
	user1 := models.User{}
	err = json.Unmarshal(byts, &user1)
	fmt.Println(err)
	assert.Equal(t, err, nil)
	fmt.Println("unmarshalled", user1)
	testPassword(t)
	err = services.GetAuthService().InternalRemoveUser("admin")
	fmt.Println(err)
	_, err = services.GetAuthService().GetUserByName("admin")
	assert.NotEqual(t, err, nil)

}

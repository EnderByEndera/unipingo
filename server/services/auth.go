package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"melodie-site/server/auth"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"net/http"

	"github.com/google/uuid"
)

type AuthService struct {
	PrivateKeys map[uuid.UUID][]byte
}

func (service *AuthService) GetAuthKey() (string, uuid.UUID) {
	prvKey, pubKey := auth.GenRsaKey()
	authUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	service.PrivateKeys[authUUID] = prvKey

	return string(pubKey), authUUID
}

func (service *AuthService) DecryptUserSecret(authUUID uuid.UUID, encryptedMessage string) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", err
	}
	if privateKey, ok := service.PrivateKeys[authUUID]; ok {
		decrypted, err := auth.RsaDecrypt(encryptedBytes, privateKey)
		if err != nil {
			return "", err
		} else {
			return string(decrypted), nil
		}
	}
	return "", errors.New("decrypt failed because authentication progress " + fmt.Sprint(authUUID) + " does not have private key.")
}

func (service *AuthService) Login(userName, password string) (user models.User, status int, err error) {
	conn := db.GetDBConn()
	user = models.User{}
	err = conn.Where("name = ?", userName).First(&user).Error

	if err != nil {
		status = http.StatusNotFound
		return
	}
	if !auth.ComparePassword(password, user.PasswordHash) {
		status = http.StatusUnauthorized
		return
	}
	status = http.StatusOK
	return
}

var authService *AuthService

func GetAuthService() *AuthService {
	if authService == nil {
		authService = &AuthService{}
		authService.PrivateKeys = map[uuid.UUID][]byte{}
	}
	return authService
}

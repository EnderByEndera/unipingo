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
	PrivateKeys       map[uuid.UUID][]byte
	WechatSessionKeys map[uuid.UUID]string
}

func (service *AuthService) GetAuthKey() (string, uuid.UUID) {
	prvKey, pubKey := auth.GenRsaKey()
	authUUID, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
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

func (service *AuthService) GetUserByWechatOpenID(openid string) (user *models.User, err error) {
	user = &models.User{}
	err = db.GetDBConn().Where("wechat_openid = ?", openid).First(user).Error
	return user, err
}

func (service *AuthService) CreateWechatUser(user *models.User) (err error) {
	err = db.GetDBConn().Create(user).Error
	return err
}

func (service *AuthService) StoreWechatSessionKey(authID uuid.UUID, sessionKey string) {
	authService.WechatSessionKeys[authID] = sessionKey
}

func (service *AuthService) GetWechatSessionKey(authID uuid.UUID) (key string, ok bool) {
	key = authService.WechatSessionKeys[authID]
	if key == "" {
		ok = false
	} else {
		ok = true
	}
	return
}

var authService *AuthService

func GetAuthService() *AuthService {
	if authService == nil {
		authService = &AuthService{}
		authService.PrivateKeys = map[uuid.UUID][]byte{}
		authService.WechatSessionKeys = map[uuid.UUID]string{}
	}
	return authService
}

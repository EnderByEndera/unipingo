package services

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"melodie-site/server/auth"
	"melodie-site/server/models"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (service *AuthService) GetUserByName(userName string) (user models.User, err error) {
	filter := bson.M{"name": userName}
	user = models.User{}
	err = getCollection("user").FindOne(context.TODO(), filter).Decode(&user)
	return
}

// 添加一位用户
// 注意不要在服务端中使用，而是在测试中用于添加用户的。
func (service *AuthService) InternalAddUser(userName, password, role string) (user models.User, err error) {
	user = models.User{Name: userName, Role: role, PasswordHash: auth.EncryptPassword(password)}

	_, err = getCollection("user").InsertOne(context.TODO(), &user)
	if err != nil {
		return
	}
	user, err = service.GetUserByName(userName)
	return
}

func (service *AuthService) InternalRemoveUser(userName string) (err error) {
	filter := bson.M{"name": userName}
	_, err = getCollection("user").DeleteOne(context.TODO(), filter)
	if err != nil {
		return
	}
	return
}

func (service *AuthService) Login(userName, password string) (user models.User, status int, err error) {
	user, err = service.GetUserByName(userName)

	if err != nil {
		status = http.StatusNotFound
		return
	}
	if !auth.ComparePassword(password, user.PasswordHash) {
		status = http.StatusBadRequest
		return
	}
	status = http.StatusOK
	return
}

func (service *AuthService) GetUserByWechatOpenID(openid string) (user *models.User, err error) {
	user = &models.User{}
	filter := bson.M{"wechatInfo.openID": openid}
	err = getCollection("user").FindOne(context.TODO(), filter).Decode(user)
	return
}

func (service *AuthService) CreateWechatUser(user *models.User) (err error) {
	_, err = getCollection("user").InsertOne(context.TODO(), user)
	if err != nil {
		return
	}
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

func (service *AuthService) GetUserByID(userID primitive.ObjectID) (user *models.User, err error) {
	filter := bson.M{"_id": userID}
	user = &models.User{}
	err = getCollection("user").FindOne(context.TODO(), filter).Decode(user)
	return
}

// PublicInfo是公开的，更新时只要有token即可，无需进行校验。
// 所以是所有PublicInfo一起更新的。
func (service *AuthService) UpdateUserPublicInfo(req *models.UserPublicInfo) (err error) {
	statement := bson.M{"$set": req}
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)
	userOID, err := primitive.ObjectIDFromHex(req.OID)
	fmt.Println(userOID)
	if err != nil {
		return
	}
	res := getCollection("user").FindOneAndUpdate(context.TODO(), bson.M{"_id": userOID}, statement, opts)
	err = res.Err()
	return
}

// PublicInfo是公开的，更新时只要有token即可，无需进行校验。
// 所以是所有PublicInfo一起更新的。
func (service *AuthService) UpdateUserSchoolInfo(userID primitive.ObjectID, schoolInfo *models.SchoolInfo) (err error) {
	statement := bson.M{"$set": bson.M{"schoolInfo": schoolInfo}}
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)
	if err != nil {
		return
	}
	res := getCollection("user").FindOneAndUpdate(context.TODO(), bson.M{"_id": userID}, statement, opts)
	err = res.Err()
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

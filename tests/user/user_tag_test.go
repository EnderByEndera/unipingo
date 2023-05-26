package user_test

import (
	"context"
	"errors"
	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"melodie-site/server/db"
	"melodie-site/server/services"
	"testing"
)

func getTags(tagNum int) (tags []string) {
	for i := 0; i < tagNum; i++ {
		tags = append(tags, primitive.NewObjectID().String())
	}
	return
}

func TestGetUserTag(t *testing.T) {
	admin, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)

	tags := make([]string, 0)
	tags, err = services.GetAuthService().GetTagsByUserID(admin.ID)
	assert.Equal(t, err, nil)

	defer func() {
		db.GetCollection("user").FindOneAndUpdate(context.TODO(), bson.M{"_id": admin.ID}, bson.M{"$set": bson.M{"userTags": tags}})
	}()

	userID := primitive.NilObjectID
	_, err = services.GetAuthService().GetTagsByUserID(userID)
	assert.Equal(t, err, errors.New("用户ID为空"))

	userID, err = primitive.ObjectIDFromHex("6464decb74edfbd3023d19d3")
	assert.Equal(t, err, nil)
	_, err = services.GetAuthService().GetTagsByUserID(userID)
	assert.Equal(t, err, errors.New("数据库查找失败"))

	genTags := getTags(5) // 随机设置为5
	db.GetCollection("user").FindOneAndUpdate(context.TODO(),
		bson.M{"_id": admin.ID},
		bson.M{"$set": bson.M{"userTags": genTags}})

	newTags, err := services.GetAuthService().GetTagsByUserID(admin.ID)

	genTagMap := make(map[string]bool)

	for _, genTag := range genTags {
		genTagMap[genTag] = true
	}

	for _, newTag := range newTags {
		assert.Equal(t, genTagMap[newTag], true)
	}
}

func TestUpdateUserTag(t *testing.T) {
	admin, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, err, nil)

	tags, err := services.GetAuthService().GetTagsByUserID(admin.ID)
	assert.Equal(t, err, nil)
	defer func() {
		db.GetCollection("user").FindOneAndUpdate(context.TODO(), bson.M{"_id": admin.ID}, bson.M{"$set": bson.M{"userTags": tags}})
	}()

	tagMap := make(map[string]bool)
	for _, tag := range tags {
		assert.Equal(t, tagMap[tag], false)
		tagMap[tag] = true
	}

	newTags := getTags(5) // 新增加的tags数量为5
	err = services.GetAuthService().UpdateUserTag(admin.ID, newTags)
	assert.Equal(t, err, nil)

	for _, tag := range newTags {
		assert.Equal(t, tagMap[tag], false)
		tagMap[tag] = true
	}

	realTags, err := services.GetAuthService().GetTagsByUserID(admin.ID)
	assert.Equal(t, err, nil)
	for _, tag := range realTags {
		assert.Equal(t, tagMap[tag], true)
	}

	userID := primitive.NilObjectID
	//assert.Equal(t, err, nil)
	err = services.GetAuthService().UpdateUserTag(userID, tags)
	assert.Equal(t, err, errors.New("用户ID为空"))

}

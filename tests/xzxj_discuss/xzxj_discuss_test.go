package xzxj_discuss_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"testing"
)

func getOneXZXJUser(user models.User) (XZXJUserForm *models.XZXJUserFormMeta) {
	XZXJUserForm = &models.XZXJUserFormMeta{
		XZXJUser: models.XZXJUser{
			ID:                primitive.NewObjectID(),
			UserID:            user.ID,
			Sections:          []string{primitive.NewObjectID().String(), primitive.NewObjectID().String()},
			Picture:           primitive.NewObjectID().String(),
			Motto:             primitive.NewObjectID().String(),
			ManagedActivities: primitive.NewObjectID().String(),
			Experience:        primitive.NewObjectID().String(),
		},
		RealName: user.RealName,
		UserTags: user.UserTags,
	}
	return
}

func TestAddXZXJUser(t *testing.T) {
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, nil, err)

	xzxjUserForm := getOneXZXJUser(user)
	userID, err := services.GetXZXJDiscussService().AddOrUpdateXZXJUser(xzxjUserForm)
	defer func() {
		_ = services.GetXZXJDiscussService().DeleteXZXJUserByID(userID)
	}()
	assert.Equal(t, nil, err)

	xzxjUser := new(models.XZXJUser)
	err = db.GetCollection("xzxj_discuss").FindOne(context.TODO(), bson.M{"userID": userID}).Decode(xzxjUser)
	assert.Equal(t, nil, err)
}

func TestDeleteXZXJUser(t *testing.T) {
	user, err := services.GetAuthService().GetUserByName("admin")
	assert.Equal(t, nil, err)

	xzxjUserForm := getOneXZXJUser(user)

	userID, err := services.GetXZXJDiscussService().AddOrUpdateXZXJUser(xzxjUserForm)
	assert.Equal(t, nil, err)

	err = services.GetXZXJDiscussService().DeleteXZXJUserByID(userID)
	assert.Equal(t, nil, err)

	err = services.GetXZXJDiscussService().DeleteXZXJUserByID(userID)
	assert.NotEqual(t, nil, err)
}

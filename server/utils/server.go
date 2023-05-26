package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"melodie-site/server/auth"
)

func SetClaims(c *gin.Context, claims auth.MelodieSiteClaims) {
	c.Set("claims", claims)
}

func GetClaims(c *gin.Context) (claims auth.MelodieSiteClaims, err error) {
	val, ok := c.Get("claims")
	if !ok {
		err = errors.New("claim does not exist")
		return
	}
	claims = val.(auth.MelodieSiteClaims)
	return
}

func GetUserID(c *gin.Context) (userID primitive.ObjectID, err error) {
	claims, err := GetClaims(c)
	if err != nil {
		return
	}
	userIDString := claims.UserID
	userID, err = primitive.ObjectIDFromHex(userIDString)
	return
}

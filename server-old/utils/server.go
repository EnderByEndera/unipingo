package utils

import (
	"errors"
	"melodie-site/server/auth"

	"github.com/gin-gonic/gin"
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

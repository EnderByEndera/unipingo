package auth

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mySigningKey = []byte("AllYourBase")

type MelodieSiteClaims struct {
	UserID string
	jwt.RegisteredClaims
}

func CreateJWTString(userID primitive.ObjectID) (string, error) {
	mm, _ := time.ParseDuration("48h")
	claims := MelodieSiteClaims{
		userID.Hex(),
		jwt.RegisteredClaims{
			// Also fixed dates can be used for the NumericDate
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(mm)),
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	return ss, err
}

func ParseJWTString(tokenString string) (*MelodieSiteClaims, bool, error) {
	if tokenString == "" || tokenString == "null" {
		return nil, false, errors.New("token string was empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &MelodieSiteClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return nil, false, err
	}
	claims, ok := token.Claims.(*MelodieSiteClaims)
	if !ok {
		return nil, token.Valid, errors.New("conversion error")
	}
	return claims, token.Valid, nil

}

func VerifyJWTString(tokenString string) error {
	claims, valid, err := ParseJWTString(tokenString)
	if err == nil && valid {
		fmt.Printf("%v", claims.UserID)
		return nil
	} else {
		fmt.Println(err)
		return err
	}
}

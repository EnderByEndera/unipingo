package routers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"melodie-site/server/auth"
	"melodie-site/server/config"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"melodie-site/server/svcerror"
	"melodie-site/server/utils"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// authProgressID: 发起临时性的加密会话的ID，一次通信完成之后，相关证书就会被销毁。

type LoginInfoRequest struct {
	AuthProgressUUID  string `json:"authProgressUUID"`
	UserName          string `json:"userName"`
	EncryptedPassword string `json:"encryptedPassword"`
}

type WechatLoginRequest struct {
	Code string `json:"code"`
}

// WechatLoginResponse 相关文档见：
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
type WechatLoginResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func CreateRSAPublicKey(c *gin.Context) {
	publicKey, authUUID := services.GetAuthService().GetAuthKey()
	c.JSON(200, map[string]string{"publicKey": publicKey, "authProgressUUID": authUUID.String()})
}

func Login(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(500, err)
	}
	req := &LoginInfoRequest{}
	err = json.Unmarshal(dataBytes, req)
	if err != nil {
		c.AbortWithError(500, err)
	}
	authProgressUUID, err := uuid.Parse(req.AuthProgressUUID)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	password, err := services.GetAuthService().DecryptUserSecret(authProgressUUID, req.EncryptedPassword)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	fmt.Println(password)

	user, status, err := services.GetAuthService().Login(req.UserName, password)

	if err != nil {
		c.String(status, err.Error())
		return
	} else {
		jwt, err := auth.CreateJWTString(user.ID, user.Role)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		} else {
			c.JSON(status, models.LoginResponse{JWTToken: jwt})
		}
	}
}

func LoginWechat(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(500, err)
	}
	reqStruct := &WechatLoginRequest{}
	err = json.Unmarshal(dataBytes, reqStruct)
	if err != nil {
		c.AbortWithError(500, err)
	}
	if reqStruct.Code == "" {
		c.String(http.StatusBadRequest, "no param 'code' input!")
		return
	}
	params := url.Values{}
	Url, err := url.Parse("https://api.weixin.qq.com/sns/jscode2session")
	if err != nil {
		return
	}
	params.Set("appid", config.GetConfig().Wechat.AppID)
	params.Set("secret", config.GetConfig().Wechat.Secret)
	params.Set("js_code", reqStruct.Code)
	params.Set("grant_type", "authorization_code")
	// params.Set("name","zhaofan")
	// params.Set("age","23")
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println(urlPath) // https://httpbin.org/get?age=23&name=zhaofan
	resp, err := http.Get(urlPath)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, errors.New("无效Wechat请求")))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	wechatLoginResponse := WechatLoginResponse{}
	err = json.Unmarshal(body, &wechatLoginResponse)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	// 如果OpenID为空，说明js_code无效
	//if wechatLoginResponse.OpenID == "" {
	//	c.Error(svcerror.New(http.StatusUnauthorized, errors.New("无效Wechat请求")))
	//	return
	//}
	user, err := services.GetAuthService().GetUserByWechatOpenID(wechatLoginResponse.OpenID)
	if err != nil {
		err := services.GetAuthService().CreateWechatUser(&models.User{
			Name: "微信用户",
			WechatInfo: models.WechatInfo{
				OpenID:  wechatLoginResponse.OpenID,
				UnionID: wechatLoginResponse.UnionID,
			},
		},
		)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		user, err = services.GetAuthService().GetUserByWechatOpenID(wechatLoginResponse.OpenID)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	authUUID, err := uuid.NewUUID()
	services.GetAuthService().StoreWechatSessionKey(authUUID, wechatLoginResponse.SessionKey)
	if err != nil {
		log.Println(err)
	}
	jwt, err := auth.CreateJWTString(user.ID, user.Role)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, models.LoginResponse{JWTToken: jwt})
}

func UploadAvatar(ctx *gin.Context) {
	f, file, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    2,
			"message": "获取数据失败",
		})
		return
	}
	defer f.Close()

	if err != nil {
		fmt.Println("获取数据失败")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "获取数据失败",
		})
	} else {
		t0 := time.Now()
		fileName, code, err := services.UploadFileByHeaderToOSS(ctx, file)
		fmt.Println(time.Since(t0))
		if err != nil {
			ctx.JSON(code, gin.H{
				"code":    1,
				"message": err.Error(),
			})
		} else {
			ctx.JSON(code, gin.H{
				"code":     0,
				"message":  "uploaded file ok!",
				"fileName": fileName,
			})
		}
	}

}

// 获取用户的公共可见信息。
// 如果没有请求参数，则直接从jwt token中读取。
func GetPublicInfo(ctx *gin.Context) {
	userIDString := ctx.Query("userID")
	if userIDString == "" {
		claims, err := utils.GetClaims(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		userIDString = claims.UserID
	}
	var userID primitive.ObjectID
	var user *models.User
	var err error
	if userIDString != "" {
		userID, err = primitive.ObjectIDFromHex(userIDString)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		user, err = services.GetAuthService().GetUserByID(userID)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, user.ToPublicInfo())
			return
		}
	} else {
		err = errors.New("userID string was empty")
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

}

// 获取用户的完整信息
func GetUser(ctx *gin.Context) {
	userIDString := ctx.Query("userID")
	if userIDString == "" {
		claims, err := utils.GetClaims(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		userIDString = claims.UserID
	}
	var userID primitive.ObjectID
	var user *models.User
	var err error
	if userIDString != "" {
		userID, err = primitive.ObjectIDFromHex(userIDString)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		user, err = services.GetAuthService().GetUserByID(userID)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, user)
			return
		}
	} else {
		err = errors.New("userID string was empty")
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
}

// 更新用户的昵称、头像等无需认证的信息
func UpdateUserPublicInfo(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(500, err)
	}
	reqStruct := &models.UserPublicInfoUpdateRequest{}
	err = json.Unmarshal(dataBytes, reqStruct)
	if err != nil {
		c.AbortWithError(500, err)
	}
	fmt.Printf("%+v\n", reqStruct)

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	err = services.GetAuthService().UpdateUserPublicInfo(userID, reqStruct)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

}

func GetUserTags(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}
	_, err = services.GetAuthService().GetUserByID(userID)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}
	tags, err := services.GetAuthService().GetTagsByUserID(userID)
	if err != nil {
		c.Error(svcerror.New(http.StatusInternalServerError, err))
		return
	}
	c.JSON(http.StatusOK, tags)
}

func UpdateUserTag(c *gin.Context) {
	dataBytes, err := c.GetRawData()
	if err != nil {
		c.Error(svcerror.New(http.StatusInternalServerError, err))
		return
	}
	req := &models.UserTagsInfoUpdateRequest{}
	err = json.Unmarshal(dataBytes, req)
	tags := req.UserTags

	if err != nil {
		c.Error(svcerror.New(http.StatusInternalServerError, err))
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	_, err = services.GetAuthService().GetUserByID(userID)
	if err != nil {
		c.Error(svcerror.New(http.StatusUnauthorized, err))
		return
	}

	err = services.GetAuthService().UpdateUserTag(userID, tags)
	if err != nil {
		c.Error(svcerror.New(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"update": true,
	})

}

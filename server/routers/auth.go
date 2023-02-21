package routers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"melodie-site/server/auth"
	"melodie-site/server/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// authProgressID: 发起临时性的加密会话的ID，一次通信完成之后，相关证书就会被销毁。

type LoginInfoRequest struct {
	AuthProgressUUID  string `json:"authProgressUUID"`
	UserName          string `json:"userName"`
	EncryptedPassword string `json:"encryptedPassword"`
}

type LoginResponse struct {
	JWTToken string `json:"jwtToken"`
}

type WechatLoginRequest struct {
	Code string `json:"code"`
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
		jwt, err := auth.CreateJWTString(user.ID)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		} else {
			c.JSON(status, LoginResponse{JWTToken: jwt})
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
	req, err := http.NewRequest("GET", "https://api.weixin.qq.com/sns/jscode2session", nil)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	// appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code

	// 返回数据示例

	q := req.URL.Query()
	q.Add("appid", "wxf7dc6cdd6711feea")
	q.Add("secret", "b4b5f723d87de6782307dda413abe99d")
	q.Add("js_code", reqStruct.Code)
	q.Add("grant_type", "authorization_code")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		c.String(http.StatusFailedDependency, err.Error())
		return
	}
	body, err := ioutil.ReadAll(req.Response.Body)
	defer response.Body.Close()
	fmt.Println(string(body))

}

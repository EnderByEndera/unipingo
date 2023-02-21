package routers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"melodie-site/server/auth"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"melodie-site/server/utils"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// 相关文档见：
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
		jwt, err := auth.CreateJWTString(user.ID)
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
	params.Set("appid", "wxf7dc6cdd6711feea")
	params.Set("secret", "b4b5f723d87de6782307dda413abe99d")
	params.Set("js_code", reqStruct.Code)
	params.Set("grant_type", "authorization_code")
	// params.Set("name","zhaofan")
	// params.Set("age","23")
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println(urlPath) // https://httpbin.org/get?age=23&name=zhaofan
	resp, err := http.Get(urlPath)
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
	user, err := services.GetAuthService().GetUserByWechatOpenID(wechatLoginResponse.OpenID)
	if err != nil {

		err := services.GetAuthService().CreateWechatUser(&models.User{
			WechatOpenID:  wechatLoginResponse.OpenID,
			WechatUnionID: wechatLoginResponse.UnionID},
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
	jwt, err := auth.CreateJWTString(user.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusAccepted, models.LoginResponse{UserInfo: *user, JWTToken: jwt})
}

func UploadAvatar(ctx *gin.Context) {
	//获取普通文本
	// name := ctx.PostForm("name")
	// 获取文件(注意这个地方的file要和html模板中的name一致)
	// file, err := ctx.FormFile("file")

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
		// fmt.Println("接收的数据", name, file.Filename)
		// //获取文件名称
		// fmt.Println(file.Filename)
		// //文件大小
		// fmt.Println(file.Size)
		// //获取文件的后缀名
		// extstring := path.Ext(file.Filename)
		// fmt.Println(extstring)
		// //根据当前时间鹾生成一个新的文件名
		// fileNameInt := time.Now().Unix()
		// fileNameStr := strconv.FormatInt(fileNameInt, 10)
		// //新的文件名
		//保存上传文件
		// filePath := filepath.Join(utils.Mkdir("upload"), "/", fileName)
		contentType := file.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		ossHandler := utils.GetOSSHandler()
		// utils.PutFile(ossHandler.Buckets.Files, file.Filename, file.Filename, contentType)
		err := utils.PutObject(ossHandler.Buckets.Files, file.Filename, contentType, f, -1)
		fmt.Println("err", err)
		// ctx.SaveUploadedFile(file, fileName)
		fmt.Println("FileName", file.Filename)
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
		})
	}

}

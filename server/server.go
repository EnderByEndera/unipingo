package server

import (
	"fmt"
	"melodie-site/server/auth"
	"melodie-site/server/routers"
	"melodie-site/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("X-Access-Token")
		fmt.Println("token", accessToken)

		// err := auth.VerifyJWTString(accessToken)
		claims, valid, err := auth.ParseJWTString(accessToken)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			c.Abort()
		}
		utils.SetClaims(c, *claims)
		if err != nil || !valid {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		//请求处理
		c.Next()
	}
}

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":8787",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}

func RunServer() {
	r := gin.Default()
	r.Use(TlsHandler())
	fileRouter := r.Group("/api/file")
	{
		fileRouter.GET("/getfile/:file", routers.SendFile)
	}
	authRouter := r.Group("/api/auth")
	{
		authRouter.POST("/rsa_public_key", routers.CreateRSAPublicKey)
		authRouter.POST("/login", routers.Login)
		authRouter.POST("/wechat_login", routers.LoginWechat)
		authRouter.POST("/upload", routers.UploadAvatar)
	}
	// articlesRouter := r.Group("/api/articles")
	// {
	// 	articlesRouter.GET("/all", authMiddleware(), routers.GetAllArticles)
	// 	articlesRouter.GET("/article", routers.GetArticle)
	// 	articlesRouter.POST("/create", authMiddleware(), routers.CreateArticle)
	// 	articlesRouter.POST("/update", authMiddleware(), routers.UpdateArticle)
	// }
	// tagsRouter := r.Group("/api/tags")
	// {
	// 	tagsRouter.GET("/tag", routers.GetTag)
	// 	tagsRouter.GET("/all", routers.GetAllTags)
	// 	tagsRouter.POST("/create", authMiddleware(), routers.CreateTag)
	// 	tagsRouter.POST("/update", authMiddleware(), routers.UpdateTag)
	// }
	// r.Run("0.0.0.0:8787")
	r.RunTLS(":8787", "cert/9325061_wechatapi.houzhanyi.com.pem", "cert/9325061_wechatapi.houzhanyi.com.key")

}

package server

import (
	"fmt"
	"melodie-site/server/auth"
	"melodie-site/server/routers"
	"melodie-site/server/services"
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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		// if origin != "" {
		// 可将将* 替换为指定的域名
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		// c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		// }

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}

func initServer() {
	_, err := services.GetAuthService().GetUserByName("admin")
	if err != nil {
		services.GetAuthService().InternalAddUser("admin", "123456")
	}
}

func RunServer() {
	r := gin.Default()
	r.Use(TlsHandler())
	r.Use(Cors())
	initServer()
	fileRouter := r.Group("/api/file")
	{
		fileRouter.GET("/getfile/:file", routers.SendFile)
	}
	authRouter := r.Group("/api/auth")
	{
		authRouter.POST("/rsaPublicKey", routers.CreateRSAPublicKey)
		authRouter.POST("/login", routers.Login)
		authRouter.POST("/wechatLogin", routers.LoginWechat)
		authRouter.POST("/upload", routers.UploadAvatar)
		authRouter.GET("/userPublicInfo", routers.GetPublicInfo)
	}
	postsRouter := r.Group("/api/posts")
	{
		postsRouter.GET("/all", routers.GetAllUserPosts)
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

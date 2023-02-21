package server

import (
	"encoding/json"
	"fmt"
	"melodie-site/server/auth"
	"melodie-site/server/db"
	"melodie-site/server/models"
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

func InitServer() {
	db.InitDB()
	svc := services.ArticleService{}
	articles, _ := svc.GetAllArticles()
	out, _ := json.MarshalIndent(articles, "", "    ")
	j := `{
        "id": 1,
        "created_at": "2022-12-31T11:31:37.579248+08:00",
        "updated_at": "2022-12-31T11:31:37.579248+08:00",
        "name": "测试",
        "abstract": "asdasdasdasda",
        "Description": "",
        "rate": 0,
        "Tags": [
            "ABM",
            "Epidemic"
        ],
        "Authors": [],
        "Links": []
    }`
	article := models.Article{}
	json.Unmarshal([]byte(j), &article)
	fmt.Printf("%+v\n", string(out))
	fmt.Println(article)
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
	authRouter := r.Group("/api/auth")
	{
		authRouter.POST("/rsa_public_key", routers.CreateRSAPublicKey)
		authRouter.POST("/login", routers.Login)
		authRouter.POST("/wechat_login", routers.LoginWechat)
	}
	articlesRouter := r.Group("/api/articles")
	{
		articlesRouter.GET("/all", authMiddleware(), routers.GetAllArticles)
		articlesRouter.GET("/article", routers.GetArticle)
		articlesRouter.POST("/create", authMiddleware(), routers.CreateArticle)
		articlesRouter.POST("/update", authMiddleware(), routers.UpdateArticle)
	}
	tagsRouter := r.Group("/api/tags")
	{
		tagsRouter.GET("/tag", routers.GetTag)
		tagsRouter.GET("/all", routers.GetAllTags)
		tagsRouter.POST("/create", authMiddleware(), routers.CreateTag)
		tagsRouter.POST("/update", authMiddleware(), routers.UpdateTag)
	}
	// r.Run("0.0.0.0:8787")
	r.RunTLS(":8787", "cert/9325061_wechatapi.houzhanyi.com.pem", "cert/9325061_wechatapi.houzhanyi.com.key")
	// handler := r.Handler()

	// http.ListenAndServeTLS("0.0.0.0:8787", "cert/9325061_wechatapi.houzhanyi.com.pem", "cert/9325061_wechatapi.houzhanyi.com.key",
	// 	handler)
	// server := http.Server{
	// 	Addr:      "0.0.0.0:8787",
	// 	Handler:   r,
	// 	TLSConfig: &tls.Config{},
	// }
	// fmt.Println(server.TLSConfig.MinVersion, server.TLSConfig.MaxVersion)
	// err := server.ListenAndServeTLS("cert/9325061_wechatapi.houzhanyi.com.pem", "cert/9325061_wechatapi.houzhanyi.com.key")
	// if err != nil {
	// 	panic(err.Error())
	// }

}

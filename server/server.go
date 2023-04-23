package server

import (
	"melodie-site/server/auth"
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
		services.GetAuthService().InternalAddUser("admin", "123456", models.RoleAdmin, nil)
	}
	_, err = services.GetAuthService().GetUserByName("demo-unpaid-user")
	if err != nil {
		services.GetAuthService().InternalAddUser("demo-unpaid-user", "123456", models.RoleUnpaidUser, nil)
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
		fileRouter.GET("/getStaticFile/:file", routers.SendStaticFile)
		fileRouter.POST("/upload", routers.UploadAvatar)
		fileRouter.POST("/uploadStaticFile", routers.UploadStaticFile)
	}
	authRouter := r.Group("/api/auth")
	{
		authRouter.POST("/rsaPublicKey", routers.CreateRSAPublicKey)
		authRouter.POST("/login", routers.Login)
		authRouter.POST("/wechatLogin", routers.LoginWechat)
		authRouter.POST("/newStuIDAuth", authMiddleware(), routers.NewStuIDAuthProc)             // 提起学生身份验证
		authRouter.GET("/unhandledStuIDAuths", authMiddleware(), routers.GetUnhandledProcs)      // 获取所有未完成的学生身份验证
		authRouter.GET("/stuIDAuth", authMiddleware(), routers.GetStuIDAuthProc)                 // 获取所有未完成的学生身份验证
		authRouter.POST("/setStuIDAuthStatus", authMiddleware(), routers.SetStudentIDAuthStatus) // 获取所有未完成的学生身份验证
		authRouter.GET("/userPublicInfo", authMiddleware(), routers.GetPublicInfo)
		authRouter.GET("/userInfo", authMiddleware(), routers.GetUser)
		authRouter.POST("/updateUserPublicInfo", authMiddleware(), routers.UpdateUserPublicInfo)
	}
	heisRouter := r.Group("/api/heis")
	{
		heisRouter.GET("/getHEIByName", routers.GetHEIByName)
		heisRouter.GET("/getHEI", routers.GetHEI)
		heisRouter.GET("/filterHEI", routers.FilterHEI)
		heisRouter.POST("/addHEIToCollection", authMiddleware(), routers.AddHEIToCollection)
		heisRouter.POST("/removeHEIFromCollection", authMiddleware(), routers.RemoveHEIFromCollection)
	}
	answersRouter := r.Group("/api/answers")
	{
		answersRouter.GET("/topics", routers.GetTopics)
		answersRouter.POST("/newAnswer", authMiddleware(), routers.NewAnswer)
		answersRouter.GET("/getAnswersRelated", authMiddleware(), routers.GetAnswersRelatedToHEIOrMajor)
	}
	majorRouter := r.Group("/api/majors")
	{
		majorRouter.GET("/getMajorByName", routers.GetMajorByName)
		majorRouter.GET("/filterMajor", routers.FilterMajor)
	}

	r.RunTLS(":8787", "cert/9325061_wechatapi.houzhanyi.com.pem", "cert/9325061_wechatapi.houzhanyi.com.key")
}

package server

import (
	"context"
	"errors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	_ "melodie-site/docs"
	"melodie-site/server/auth"
	"melodie-site/server/config"
	"melodie-site/server/models"
	"melodie-site/server/routers"
	"melodie-site/server/services"
	"melodie-site/server/svcerror"
	"melodie-site/server/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("X-Access-Token")

		// err := auth.VerifyJWTString(accessToken)
		claims, valid, err := auth.ParseJWTString(accessToken)
		if err != nil {
			c.Error(svcerror.New(http.StatusBadRequest, err))
			c.Abort()
			return
		}
		utils.SetClaims(c, *claims)
		if err != nil || !valid {
			c.Error(svcerror.New(http.StatusUnauthorized, err))
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
			c.Abort()
			return
		}

		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if e := c.Errors.Last(); e != nil {
			err := e.Err
			if svcErr, ok := err.(*svcerror.SvcErr); ok {
				c.JSON(svcErr.Code, svcErr)
			} else {
				unexpectedErr := svcerror.AppendData(svcerror.New(http.StatusInternalServerError, errors.New("internal server error")), err)
				c.JSON(http.StatusInternalServerError, unexpectedErr)
			}
		}
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
		_, err = services.GetAuthService().InternalAddUser("admin", "123456", models.RoleAdmin, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = services.GetAuthService().GetUserByName("demo-unpaid-user")
	if err != nil {
		_, err = services.GetAuthService().InternalAddUser("demo-unpaid-user", "123456", models.RoleUnpaidUser, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func endServer() (err error) {
	err = services.GetAuthService().InternalRemoveUser("admin")
	if err != nil {
		return
	}

	err = services.GetAuthService().InternalRemoveUser("demo-unpaid-user")
	if err != nil {
		return
	}
	return
}

func RunServer() {
	r := gin.Default()
	r.Use(TlsHandler())
	r.Use(Cors())
	r.Use(ErrorHandler())
	initServer()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	serviceRouter := r.Group("/api")
	{
		fileRouter := serviceRouter.Group("/file")
		{
			fileRouter.GET("/get/:file", routers.SendFile)
			fileRouter.GET("/getStaticFile/:file", routers.SendStaticFile)
			fileRouter.POST("/upload", routers.UploadAvatar)
			fileRouter.POST("/uploadStaticFile", routers.UploadStaticFile)
		}
		authRouter := serviceRouter.Group("/auth")
		{
			authRouter.POST("/rsaPublicKey", routers.CreateRSAPublicKey)
			authRouter.POST("/login", routers.Login)
			authRouter.POST("/wechatLogin", routers.LoginWechat)
			authRouter.POST("/newStuIDAuth", authMiddleware(), routers.NewStuIDAuthProc)             // 提起学生身份验证
			authRouter.POST("/updateStuIDAuth", authMiddleware(), routers.UpdateStuIDAuthProc)       // 修改学生身份验证信息
			authRouter.GET("/unhandledStuIDAuths", authMiddleware(), routers.GetUnhandledProcs)      // 获取所有未完成的学生身份验证
			authRouter.GET("/stuIDAuth", authMiddleware(), routers.GetStuIDAuthProc)                 // 获取一个学生身份验证的实例
			authRouter.POST("/setStuIDAuthStatus", authMiddleware(), routers.SetStudentIDAuthStatus) // 设置学生身份验证的状态
			authRouter.GET("/userPublicInfo", authMiddleware(), routers.GetPublicInfo)
			authRouter.GET("/userInfo", authMiddleware(), routers.GetUser)
			authRouter.POST("/updateUserPublicInfo", authMiddleware(), routers.UpdateUserPublicInfo)
			authTagRouter := authRouter.Group("/tag")
			{
				authTagRouter.GET("/get", authMiddleware(), routers.GetUserTags)
				authTagRouter.POST("/update", authMiddleware(), routers.UpdateUserTag)
			}
		}
		heisRouter := serviceRouter.Group("/heis")
		{
			heisRouter.GET("/getHEIByName", routers.GetHEIByName)
			heisRouter.GET("/getHEI", routers.GetHEI)
			heisRouter.GET("/filterHEI", routers.FilterHEI)
			heisRouter.POST("/addHEIToCollection", authMiddleware(), routers.AddHEIToCollection)
			heisRouter.POST("/removeHEIFromCollection", authMiddleware(), routers.RemoveHEIFromCollection)
		}
		answersRouter := serviceRouter.Group("/answers")
		{
			answersRouter.GET("/topics", routers.GetTopics)
			answersRouter.POST("/newAnswer", authMiddleware(), routers.NewAnswer)
			answersRouter.GET("/getAnswersRelated", authMiddleware(), routers.GetAnswersRelatedToHEIOrMajor)
			answersRouter.POST("/approveOrDisapprove", authMiddleware(), routers.ApproveOrDisapproveAnswer)
		}
		majorRouter := r.Group("/majors")
		{
			majorRouter.GET("/getMajorByName", routers.GetMajorByName)
			majorRouter.GET("/filterMajor", routers.FilterMajor)
		}
		orderRouter := serviceRouter.Group("/orders")
		{
			// TODO: 具体API命名还需要和前端商定
			orderRouter.POST("/prepay", authMiddleware(), routers.PrepayOrder)
			orderRouter.POST("/notify", authMiddleware(), routers.NotifyOrder)
			orderRouter.POST("/status/get", authMiddleware(), routers.GetOrderStatus)
			orderRouter.POST("/cancel", authMiddleware(), routers.CancelOrder)
		}
		questionBoxRouter := serviceRouter.Group("/questionbox")
		{
			qbQuestionRouter := questionBoxRouter.Group("/question")
			{
				qbQuestionRouter.POST("/new", authMiddleware(), routers.NewQuestion)
				qbQuestionRouter.GET("/query", authMiddleware(), routers.QueryQuestionByID)
				qbQuestionRouter.GET("/list", authMiddleware(), routers.QueryMyQuestionList)
				qbQuestionRouter.POST("/description/update", authMiddleware(), routers.UpdateQuestionDescription)
				qbQuestionRouter.POST("/school/update", authMiddleware(), routers.UpdateQuestionSchoolOrMajor)
				qbQuestionRouter.POST("/major/update", authMiddleware(), routers.UpdateQuestionSchoolOrMajor)
			}
			qbLabelRouter := questionBoxRouter.Group("/label")
			{
				qbLabelRouter.POST("/new", authMiddleware(), routers.NewLabels)
				qbLabelRouter.GET("/user/get", authMiddleware(), routers.GetLabelsFromUser)
				qbLabelRouter.POST("/question/get", authMiddleware(), routers.GetLabelsFromQuestion)
				qbLabelRouter.POST("/delete", authMiddleware(), routers.DeleteLabel)
				qbLabelRouter.POST("/content/update", authMiddleware(), routers.UpdateLabelContent)
			}
			qbAnswerRouter := questionBoxRouter.Group("/answer")
			{
				qbAnswerRouter.POST("/new", authMiddleware(), routers.NewQuestionBoxAnswer)
				qbAnswerRouter.GET("/query", authMiddleware(), routers.QueryAnswerByID)
				qbAnswerRouter.GET("/list", authMiddleware(), routers.GetAnswerList)
				qbAnswerRouter.GET("/mylist", authMiddleware(), routers.GetMyAnswerList)
				qbAnswerRouter.POST("/content/update", authMiddleware(), routers.UpdateAnswerContent)
				qbAnswerRouter.POST("/read", authMiddleware(), routers.ReadAnswerByUser)
			}
		}
		xzxjDiscussRouter := serviceRouter.Group("/xzxjdiscuss")
		{
			xzxjDiscussRouter.POST("/add", authMiddleware(), routers.AddOrUpdateXZXJUser)
			xzxjDiscussRouter.POST("/update", authMiddleware(), routers.AddOrUpdateXZXJUser)
			xzxjDiscussRouter.GET("/query", authMiddleware(), routers.QueryXZXJUserByUserID)
			xzxjDiscussRouter.GET("/delete", authMiddleware(), routers.DeleteXZXJUser)
		}
	}

	srv := &http.Server{
		Addr:    ":8787",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServeTLS(
			config.GetConfig().Server.CertFile,
			config.GetConfig().Server.KeyFile); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // Ctrl-C & kill
	<-quit

	log.Println("Server Shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), config.GetConfig().Server.Timeout)
	defer cancel()
	//if err := endServer(); err != nil {
	//	log.Println("Server failed to delete users admin and demo-unpaid-user")
	//}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server failed to shutdown. Err: ", err)
	}
	log.Println("Server Exited")
}

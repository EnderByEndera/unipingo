# 大学声 服务端
## 简介

目前服务端既是“开发环境”也是“测试环境”。这是因为我们的开发环境是单体应用，开发和测试环境没有太多区别。如果自行配置环境反而可能有些麻烦。

未来正式部署时，将部署专门的生产环境。

地址：
- 81.70.235.52

密码：
- 私发在开发群

建议使用Visual Studio Code中的SSH远程开发工具进行开发。


## 启动项目
```sh
go run main.go run
```
会看到如下的输出。输出中显示了所有的服务端路由及其方法。

注意，现在有关用户管理的内容还没有写好，还需要进一步的工作。
```sh
[hzy@VM-0-8-centos internet-plus-backend]$ go run main.go run
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /api/file/getfile/:file   --> melodie-site/server/routers.SendFile (4 handlers)
[GIN-debug] POST   /api/auth/rsa_public_key  --> melodie-site/server/routers.CreateRSAPublicKey (4 handlers)
[GIN-debug] POST   /api/auth/login           --> melodie-site/server/routers.Login (4 handlers)
[GIN-debug] POST   /api/auth/wechat_login    --> melodie-site/server/routers.LoginWechat (4 handlers)
[GIN-debug] POST   /api/auth/upload          --> melodie-site/server/routers.UploadAvatar (4 handlers)
[GIN-debug] GET    /api/articles/all         --> melodie-site/server/routers.GetAllArticles (5 handlers)
[GIN-debug] GET    /api/articles/article     --> melodie-site/server/routers.GetArticle (4 handlers)
[GIN-debug] POST   /api/articles/create      --> melodie-site/server/routers.CreateArticle (5 handlers)
[GIN-debug] POST   /api/articles/update      --> melodie-site/server/routers.UpdateArticle (5 handlers)
[GIN-debug] GET    /api/tags/tag             --> melodie-site/server/routers.GetTag (4 handlers)
[GIN-debug] GET    /api/tags/all             --> melodie-site/server/routers.GetAllTags (4 handlers)
[GIN-debug] POST   /api/tags/create          --> melodie-site/server/routers.CreateTag (5 handlers)
[GIN-debug] POST   /api/tags/update          --> melodie-site/server/routers.UpdateTag (5 handlers)
[GIN-debug] Listening and serving HTTPS on :8787
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
```
## 文件结构
主要位于/server目录下，由根目录下的main.go文件引用。目录结构如下：
```
server/
|-- auth    // 用户认证的基础方法，如jwt\rsa\sha256等。
|   |-- jwt.go
|   |-- rsa.go
|   `-- sha256.go
|-- config  // 用于设置读取的文件。
|   `-- config.go
|-- db      // 数据库连接
|   |-- db.go
|   `-- mongo.go
|-- models  // 数据模型
|   |-- article.go
|   |-- auth.go
|   |-- base.go
|   |-- post.go
|   `-- tags.go
|-- routers // 路由文件，定义路由函数。
|   |-- article.go
|   |-- auth.go
|   |-- oss.go
|   `-- tag.go
|-- server.go// 服务端程序定义的文件
|-- services // 服务层文件，定义服务函数，由路由函数调用，负责数据库连接、查询、修改，权限认证等。
|   |-- article.go
|   |-- auth.go
|   |-- oss.go
|   |-- posts.go
|   `-- tags.go
`-- utils   // 乱七八糟的服务函数。
    |-- oss.go
    `-- server.go
```
package main

import (
	"melodie-site/server"
)

// @title UniPingo Backend
// @version 0.0.1
// @description This is the backend for UniPingo application
// @termsOfService http://swagger.io/terms

// @contact.name Songyue Chen
// @contact.email enderbybear@foxmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1
// @BasePath /api
// @query.collection.format multi
func main() {
	server.RunServer()
}

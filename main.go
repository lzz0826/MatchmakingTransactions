package main

import (
	"TradeMatching/common/mysql"
	"TradeMatching/config"
	_ "TradeMatching/docs"
	"TradeMatching/middleware"
	routes "TradeMatching/route"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

var HttpServer *gin.Engine

// @title TradeMatching
// @version 1.0
// @description 股票撮合系統
// @contact.name tony
// @contact.url
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8081
// schemes http
func main() {
	// 服务停止时清理数据库链接
	defer func() {
		if mysql.GormDb != nil {
			sqlDB, err := mysql.GormDb.DB()
			if err != nil {
				panic("failed to get DB instance: " + err.Error())
			}
			_ = sqlDB.Close()
		}
	}()
	//配置文件
	config.InitViper()

	//初始化mysql
	mysql.InitMysql()

	//HTTP 启动服务
	RunHttp()

}

// RunHttp 配置并启动服务
func RunHttp() {
	// 服务配置
	serverConfig := config.GetServerConfig()

	// gin 运行时 release debug test
	gin.SetMode(serverConfig["ENV"])

	HttpServer = gin.Default()

	//使用自訂上下文
	HttpServer.Use(middleware.TraceMiddleware())

	// Swagger UI
	HttpServer.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册路由
	routes.RegisterRoutes(HttpServer)

	serverAddr := serverConfig["HOST"] + ":" + serverConfig["PORT"]

	// 启动服务
	err := HttpServer.Run(serverAddr)

	if nil != err {
		panic("config run errors: " + err.Error())
	}
}

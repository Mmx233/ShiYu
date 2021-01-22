package Router

import (
	"Mmx/Middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	//Gin初始化
	G := gin.Default()

	//启动防洪协程
	go Middlewares.Sec.InitIpLogger()

	//中间件使用
	G.Use(Middlewares.Sec.Main)  //安全中间件
	G.Use(Middlewares.Auth.Main) //鉴权中间件

	//路由分级
	E := G.Group("/api")

	//分组
	routerAdmin(E.Group("/admin")) //管理员
	routerUser(E.Group("/user"))   //用户
	routerBiz(E.Group("/biz"))     //商户
	routerPublic(E.Group("/v3"))   //公共接口

	G.Run(":1986")
}

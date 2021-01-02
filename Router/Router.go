package Router

import (
	"Mmx/Middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouter(){
	//Gin初始化
	G:=gin.Default()

	//启动防洪协程
	Middlewares.Sec.InitIpLogger()

	//中间件使用
	G.Use(Middlewares.Sec.Main)//安全中间件
	G.Use(Middlewares.Auth.Main)//鉴权中间件

	//路由分组
	routerAdmin(G.Group("/api/admin"))//管理员
	routerUser(G.Group("/api/user"))//用户
	routerBiz(G.Group("/api/biz"))//商户
}

package Router

import (
	"Mmx/Middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouter(){
	G:=gin.Default()//Gin初始化

	Middlewares.Sec.InitIpLogger()//启动防洪协程

	//中间件使用
	G.Use(Middlewares.Sec.Main)//安全中间件
	G.Use(Middlewares.Auth.Main)//鉴权中间件
	//路由分组
	////路由分组
	/*admin:=G.Group("/api/admin")//管理员
	user:=G.Group("/api/user")//用户
	biz:=G.Group("/api/biz")//商户
		menu:=biz.Group("/menu")
		cat:=biz.Group("/cat")
		search:=biz.Group("/search")*/


}

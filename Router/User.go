package Router

import (
	Controllers "Mmx/Controllers/User"
	"Mmx/Middlewares"
	"github.com/gin-gonic/gin"
)

func routerUser(G *gin.RouterGroup) {
	G.Use(Middlewares.Check.UserName)
	G.GET("/:username", Controllers.User.Information) //获取账户信息
	G.PUT("/:username", Controllers.User.Renew)       //修改信息
	G.PATCH("/:username", Controllers.User.Change)    //修改指定属性
	G.DELETE("/:username", Controllers.User.Delete)   //删除用户
}

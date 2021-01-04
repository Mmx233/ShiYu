package Router

import (
	Controllers "Mmx/Controllers/User"
	"Mmx/Middlewares"
	"github.com/gin-gonic/gin"
)

func routerUser(G *gin.RouterGroup){
	G.Use(Middlewares.Check.UserName)
	G.GET("/:username",Controllers.User.Information)//获取账户信息
}

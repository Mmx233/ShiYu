package Router

import (
	"Mmx/Controllers"
	"github.com/gin-gonic/gin"
)

func routerPublic(G *gin.RouterGroup){
	G.POST("/login",Controllers.Public.Login)//登录
	G.POST("/register",Controllers.Public.Register)//注册
}

package Router

import (
	Controllers "Mmx/Controllers/Admin"
	"Mmx/Middlewares"
	"github.com/gin-gonic/gin"
)

func routerAdmin(G *gin.RouterGroup){
	G.Use(Middlewares.Check.UserName)
	G.POST("/:username",Controllers.Admin.NewAccount)//新账户
	G.GET("/:username",Controllers.Admin.Information)//获取信息
	G.PUT("/:username",Controllers.Admin.Renew)//修改信息
	G.PATCH("/:username",Controllers.Admin.Change)//修改指定信息
	G.DELETE("/:username",Controllers.Admin.Delete)//删除账号
}
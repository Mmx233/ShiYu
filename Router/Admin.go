package Router

import (
	Controllers "Mmx/Controllers/Admin"
	"Mmx/Middlewares"
	"github.com/gin-gonic/gin"
)

func routerAdmin(G *gin.RouterGroup){
	G.Use(Middlewares.Check.UserName)
	G.POST("/:username",Controllers.Admin.NewAccount)//新账户
}

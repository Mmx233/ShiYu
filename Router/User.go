package Router

import (
	"Mmx/Middlewares"
	"github.com/gin-gonic/gin"
)

func routerUser(G *gin.RouterGroup){
	G.Use(Middlewares.Check.UserName)

}

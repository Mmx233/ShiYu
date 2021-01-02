package Router

import "github.com/gin-gonic/gin"

func routerBiz(G *gin.RouterGroup){
	menu:=G.Group("/menu")
	cat:=G.Group("/cat")
	search:=G.Group("/search")
}
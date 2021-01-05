package Router

import (
	Controllers "Mmx/Controllers/Biz"
	"github.com/gin-gonic/gin"
)

func routerBiz(G *gin.RouterGroup) {
	//商户
	G.GET("/list", Controllers.Biz.ListBiz) //列出商户

	//菜单
	//menu:=G.Group("/menu")

	//分类
	//cat:=G.Group("/cat")

	//搜索
	//search:=G.Group("/search")
}

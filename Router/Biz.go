package Router

import (
	Controllers "Mmx/Controllers/Biz"
	"github.com/gin-gonic/gin"
)

func routerBiz(G *gin.RouterGroup) {
	//商户
	G.GET("/list", Controllers.Biz.ListBiz) //列出商户
	G.GET("/", Controllers.Biz.Information) //获取商户信息
	G.POST("/",Controllers.Biz.New)//新增商户
	G.PUT("/",Controllers.Biz.Renew)//修改商户
	G.PATCH("/",Controllers.Biz.Change)//修改商户特定属性
	G.DELETE("/",Controllers.Biz.Delete)//删除商家

	//菜单
	//menu:=G.Group("/menu")

	//分类
	//cat:=G.Group("/cat")

	//搜索
	//search:=G.Group("/search")
}

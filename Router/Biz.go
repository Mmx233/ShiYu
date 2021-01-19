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
	menu:=G.Group("/menu")
	menu.GET("/",Controllers.Biz.Menu.Information)//根据商家获取菜单
	menu.GET("/cat")//根据分类获取菜单


	//分类
	cat:=G.Group("/cat")
	cat.GET("/list",Controllers.Biz.Cat.ListContent)//列出分类下商户/菜品
	cat.GET("/",Controllers.Biz.Cat.List)//列出分类
	cat.POST("/",Controllers.Biz.Cat.New)//新增分类
	cat.PUT("/",Controllers.Biz.Cat.Renew)//修改分类
	cat.DELETE("/",Controllers.Biz.Cat.Delete)//删除分类

	//搜索
	//search:=G.Group("/search")
}

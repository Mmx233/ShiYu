package Router

import (
	Controllers "Mmx/Controllers/Biz"
	"github.com/gin-gonic/gin"
)

func routerBiz(G *gin.RouterGroup) {
	//商户
	G.GET("/list", Controllers.Biz.ListBiz) //列出商户
	G.GET("/", Controllers.Biz.Information) //获取商户信息
	G.POST("/", Controllers.Biz.New)        //新增商户
	G.PUT("/", Controllers.Biz.Renew)       //修改商户
	G.PATCH("/", Controllers.Biz.Change)    //修改商户特定属性
	G.DELETE("/", Controllers.Biz.Delete)   //删除商家

	//以下均为Biz子路由

	//菜单
	menu := G.Group("/menu")
	menu.GET("/", Controllers.Biz.Menu.Information)                      //根据商家获取菜单
	menu.GET("/cat", Controllers.Biz.Menu.InformationForCat)             //根据分类获取菜单
	menu.GET("/favorites", Controllers.Biz.Menu.InformationForFavorites) //获取favorite的菜单
	menu.POST("/", Controllers.Biz.Menu.New)                             //新建菜单
	menu.PUT("/", Controllers.Biz.Menu.Renew)                            //修改菜单
	menu.PATCH("/", Controllers.Biz.Menu.Change)                         //修改菜单部分
	menu.DELETE("/", Controllers.Biz.Menu.Delete)                        //删除菜单

	//分类
	cat := G.Group("/cat")
	cat.GET("/list", Controllers.Biz.Cat.ListContent) //列出分类下商户
	cat.GET("/", Controllers.Biz.Cat.List)            //列出分类
	cat.POST("/", Controllers.Biz.Cat.New)            //新增分类
	cat.PUT("/", Controllers.Biz.Cat.Renew)           //修改分类
	cat.DELETE("/", Controllers.Biz.Cat.Delete)       //删除分类

	//favorite
	fav:=menu.Group("/fav")
	fav.POST("/",Controllers.Biz.Menu.Fav.Make)//收藏
	fav.DELETE("/",Controllers.Biz.Menu.Fav.Cancel)//取消收藏

	//搜索
	search:=G.Group("/search")
	search.GET("/biz",Controllers.Biz.Search.Biz)//搜索商家
	search.GET("/menu",Controllers.Biz.Search.Menu)//搜索菜单
	search.GET("/fav",Controllers.Biz.Search.Fav)//搜索收藏
}

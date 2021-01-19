package Controllers

import (
	"Mmx/Modules"
	"Mmx/Service"
	"github.com/gin-gonic/gin"
)

type menu struct {}

func (*menu)Information(c *gin.Context){
	type FORM struct {
		Id uint `form:"id" binding:"required,min=1"`
	}
	var form FORM
	if !Modules.Tool.BindForm(c,&form){
		return
	}
	if Service.Checker.BizExist(form.Id){
		Modules.CallBack.Error(c,116)
		return
	}
	type d struct{
		Id uint `json:"id"`
		BizId uint `json:"biz_id"`
		Name string `json:"name"`
		Price string `json:"price"`
		CatFoodId uint `json:"cat_food_id"`
		IsFavorite bool `json:"is_favorite" skip:"true"`
	}
	data:=make([]d,1)
	if Service.Get(c,"biz_menu",&data, map[string]interface{}{
		"biz_id":form.Id,
	})!=nil{
		return
	}
	//favorite相关
	var u struct{
		Favorites []uint `json:"favorites"`
	}
	if Service.GetRow(c,"user",&u, map[string]interface{}{
		"username":c.Get("username"),
	})!=nil{
		return
	}
	for i:=0;i<len(data);i++{
		//Go没有find，好烦啊
		data[i].IsFavorite=Modules.Tool.Find(u.Favorites,data[i].Id)
	}
	Modules.CallBack.Success(c,data)
}


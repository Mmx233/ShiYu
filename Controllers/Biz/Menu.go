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
	var data struct{
		Id uint `json:"id"`
		BizId uint `json:"biz_id"`
		Name string `json:"name"`
		Price string `json:"price"`
		CatFoodId uint `json:"cat_food_id"`
		IsFavorite string `json:"is_favorite" skip:"true"`
	}
}
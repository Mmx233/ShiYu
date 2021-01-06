package Controllers

import (
	"Mmx/Modules"
	"github.com/gin-gonic/gin"
)

type biz struct{}

var Biz biz

type bizData struct {
	Id uint `json:"id"`
	Name string `json:"name"`
	PicUrl []uint `json:"pic_url"`
	Grade float32 `json:"grade"`
	GradeWeight uint `json:"grade_weight"`
	MaxPrice float32 `json:"max_price"`
	MinPrice float32 `json:"min_price"`
	Address string `json:"address"`
	CatAddress string `json:"cat_address"`
	Husk uint32 `json:"husk"`
	Share uint32 `json:"share"`
	Favorite uint `json:"favorite"`
	Dislike uint `json:"dislike"`
}

func (*biz) ListBiz(c *gin.Context) {
	type listForm struct {
		Limit uint `form:"limit" binding:"required,min=1"`
		Page uint `form:"page" binding:"required,min=1"`
	}
	var form listForm
	if !Modules.Tool.BindForm(c,&form){
		return
	}

}

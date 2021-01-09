package Controllers

import (
	"Mmx/Modules"
	"Mmx/Service"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type biz struct{}

var Biz biz

type bizData struct {
	Id uint `json:"id"`
	Name string `json:"name"`
	PicUrl []string `json:"pic_url"`
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
	var data =make([]bizData,form.Limit)
	if Service.Get(c,"biz",&data, map[string]interface{}{},int(form.Page))!=nil{
		return
	}
	for i:=range data{//处理img Url
		Modules.Tool.ImgString(c,"biz",&data[i].PicUrl)
	}
	Modules.CallBack.Success(c,data)
}

func (*biz)Information(c *gin.Context){
	type infoForm struct {
		Id uint `form:"id" binding:"required,min=1"`
	}
	var form infoForm
	if !Modules.Tool.BindForm(c,&form){
		return
	}
	if !Service.Checker.BizExist(form.Id){
		Modules.CallBack.Error(c,116)
		return
	}
	var data bizData
	if Service.GetRow(c,"biz",&data, map[string]interface{}{
		"id":form.Id,
	})!=nil{
		return
	}
	Modules.Tool.ImgString(c,"biz",&data.PicUrl)//处理img url
	Modules.CallBack.Success(c,data)
}

func (*biz)New(c *gin.Context){
	type newBizForm struct {
		Name string `form:"name" binding:"required"`
		Pic []string `form:"pic" binding:"required"`
		Address string `form:"address" binding:"required"`
		CatAddressId uint `form:"cat_address_id" binding:"required"`
	}
	var form newBizForm
	if !Modules.Tool.BindForm(c,&form){
		return
	}
	if !Modules.Checker.Form(c,&form){
		return
	}
	if !Service.Checker.CatIdExist(1,form.CatAddressId){
		Modules.CallBack.Error(c,119)
		return
	}
	picString,_:=json.Marshal(form.Pic)
	if _,err:=Service.Insert(c,"biz", map[string]interface{}{
		"name":form.Name,
		"pic_url":picString,
		"grade":0,
		"grade_weight":0,
		"max_price":0,
		"min_price":0,
		"address":form.Address,
		"husk":0,
		"share":0,
		"favorite":0,
		"dislike":0,
		"cat_address_id":form.CatAddressId,
	});err!=nil{
		return
	}
	Modules.CallBack.Default(c)
}
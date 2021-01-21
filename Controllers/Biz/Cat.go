package Controllers

import (
	"Mmx/Modules"
	"Mmx/Service"
	"github.com/gin-gonic/gin"
)

type cat struct{}

func (*cat) ListContent(c *gin.Context) {
	var form struct {
		Cat   string `form:"cat" binding:"required"`
		Id    uint   `form:"id" binding:"required"`
		Limit uint   `form:"limit" binding:"required,min=1"`
		Page  uint   `form:"page" binding:"required,min=1"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Modules.Checker.Form(c, &form) {
		return
	}
	if !Service.Checker.CatIdExist(form.Cat, form.Id) {
		Modules.CallBack.Error(c, 119)
		return
	}
	var data interface{}
	if form.Cat == "address" {
		data = make([]bizData, form.Limit)
		if Service.GetWithLimit(c, "biz", &data, map[string]interface{}{
			"cat_address_id": form.Id,
		}, int(form.Page)) != nil {
			return
		}
	} else {
		//DEMO 菜单获取
	}
	Modules.CallBack.Success(c, &data)
}

func (*cat) List(c *gin.Context) {
	var form struct {
		Cat string `form:"cat" binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Modules.Checker.Form(c, &form) {
		return
	}
	type d struct {
		Id   uint   `json:"id"`
		Name string `json:"name"`
	}
	data := make([]d, 1)
	if Service.Get(c, "cat_"+form.Cat, &data, map[string]interface{}{}) != nil {
		return
	}
	Modules.CallBack.Success(c, data)
}

func (*cat) New(c *gin.Context) {
	var form struct {
		Cat  string `form:"cat" binding:"required"`
		Name string `form:"name" binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Modules.Checker.Form(c, &form) {
		return
	}
	id, err := Service.Insert(c, "cat_"+form.Cat, map[string]interface{}{
		"name": form.Name,
	})
	if err != nil {
		return
	}
	Modules.CallBack.Success(c, map[string]int64{
		"id": id,
	})
}

func (*cat) Renew(c *gin.Context) {
	var form struct {
		Cat  string `form:"cat" binding:"required"`
		Id   uint   `form:"id" binding:"required"`
		Name string `form:"name" binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Modules.Checker.Form(c, &form) {
		return
	}
	if !Service.Checker.CatIdExist(form.Cat, form.Id) {
		Modules.CallBack.Error(c, 119)
		return
	}
	if _, err := Service.Update(c, "cat_"+form.Cat, map[string]interface{}{
		"name": form.Name,
	}, map[string]interface{}{
		"id": form.Id,
	}); err != nil {
		return
	}
	Modules.CallBack.Default(c)
}

func (*cat) Delete(c *gin.Context) {
	var form struct {
		Cat string `form:"cat" binding:"required"`
		Id  uint   `form:"id" binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Modules.Checker.Form(c, &form) {
		return
	}
	if !Service.Checker.CatIdExist(form.Cat, form.Id) {
		Modules.CallBack.Error(c, 119)
		return
	}
	if _, err := Service.Delete(c, "cat_"+form.Cat, map[string]interface{}{
		"id": form.Id,
	}); err != nil {
		return
	}
	Modules.CallBack.Default(c)
}

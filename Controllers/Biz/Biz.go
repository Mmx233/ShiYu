package Controllers

import (
	"Mmx/Modules"
	"Mmx/Service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strings"
)

type biz struct {
	Cat    cat
	Menu   menu
	Search sear
}

var Biz biz

type bizData struct {
	Id           uint     `json:"id"`
	Name         string   `json:"name"`
	PicUrl       []string `json:"pic_url"`
	Grade        float32  `json:"grade"`
	GradeWeight  uint     `json:"grade_weight"`
	MaxPrice     float32  `json:"max_price"`
	MinPrice     float32  `json:"min_price"`
	Address      string   `json:"address"`
	CatAddressId uint     `json:"cat_address_id"`
	Husk         uint32   `json:"husk"`
	Share        uint32   `json:"share"`
	Favorite     uint     `json:"favorite"`
}

func (*biz) ListBiz(c *gin.Context) {
	var form struct {
		Limit uint `form:"limit" binding:"required,min=1"`
		Page  uint `form:"page" binding:"required,min=1"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	var data = make([]bizData, form.Limit)
	if Service.GetWithLimit(c, "biz", &data, nil, int(form.Page)) != nil {
		return
	}
	for i := range data { //处理img Url
		Modules.Tool.ImgString(c, "biz", &data[i].PicUrl)
	}
	Modules.CallBack.Success(c, data)
}

func (*biz) Information(c *gin.Context) {
	var form struct {
		Id uint `form:"id" binding:"required,min=1"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Service.Checker.BizExist(form.Id) {
		Modules.CallBack.Error(c, 116)
		return
	}
	var data bizData
	if Service.GetRow(c, "biz", &data, map[string]interface{}{
		"id": form.Id,
	}) != nil {
		return
	}
	Modules.Tool.ImgString(c, "biz", &data.PicUrl) //处理img url
	Modules.CallBack.Success(c, data)
}

func (*biz) New(c *gin.Context) {
	var form struct {
		Name         string   `form:"name" binding:"required"`
		Pic          []string `form:"pic"`
		Address      string   `form:"address" binding:"required"`
		CatAddressId uint     `form:"cat_address_id" binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Modules.Checker.Form(c, &form) {
		return
	}
	if !Service.Checker.CatIdExist("address", form.CatAddressId) {
		Modules.CallBack.Error(c, 119)
		return
	}
	if Service.Checker.NameExist("biz", form.Name, nil) {
		Modules.CallBack.Error(c, 120)
		return
	}
	id, err := Service.Insert(c, "biz", map[string]interface{}{
		"name":           form.Name,
		"pic_url":        form.Pic,
		"grade":          0,
		"grade_weight":   0,
		"max_price":      0,
		"min_price":      0,
		"address":        form.Address,
		"husk":           0,
		"share":          0,
		"favorite":       0,
		"cat_address_id": form.CatAddressId,
	})
	if err != nil {
		return
	}
	Modules.CallBack.Success(c, map[string]int64{
		"id": id,
	})
}

func (*biz) Renew(c *gin.Context) {
	var form struct {
		Id           uint     `form:"id" binding:"required"`
		Name         string   `form:"name" binding:"required"`
		Pic          []string `form:"pic"`
		Grade        uint     `form:"grade"`
		GradeWeight  uint     `form:"grade_weight"`
		Address      string   `form:"address" binding:"required"`
		CatAddressId uint     `form:"cat_address_id" binding:"required"`
		Husk         uint     `form:"husk"`
		Share        uint     `form:"share"`
		Favorite     uint     `form:"favorite"`
		MaxPrice     uint     `form:"max_price"`
		MinPrice     uint     `form:"min_price"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Modules.Checker.Form(c, &form) {
		return
	}
	if !Service.Checker.BizExist(form.Id) {
		Modules.CallBack.Error(c, 116)
		return
	}
	if !Service.Checker.CatIdExist("address", form.CatAddressId) {
		Modules.CallBack.Error(c, 119)
		return
	}
	if Service.Checker.NameExist("biz", form.Name, nil) {
		Modules.CallBack.Error(c, 120)
		return
	}
	if _, err := Service.Update(c, "biz", map[string]interface{}{
		"name":           form.Name,
		"pic_url":        form.Pic,
		"grade":          form.Grade,
		"grade_weight":   form.GradeWeight,
		"max_price":      form.MaxPrice,
		"min_price":      form.MinPrice,
		"address":        form.Address,
		"husk":           form.Husk,
		"share":          form.Share,
		"favorite":       form.Favorite,
		"cat_address_id": form.CatAddressId,
	}, map[string]interface{}{
		"id": form.Id,
	}); err != nil {
		return
	}
	Modules.CallBack.Default(c)
}

func (*biz) Change(c *gin.Context) {
	var form struct {
		Id     uint        `form:"id" binding:"required"`
		Target string      `form:"target" binding:"required"`
		Value  interface{} `form:"value"  binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Service.Checker.BizExist(form.Id) {
		Modules.CallBack.Error(c, 116)
		return
	}
	form.Target = strings.ToLower(form.Target)
	switch form.Target {
	case "name":
		if _, ok := form.Value.(string); !ok {
			Modules.CallBack.Error(c, 101)
			return
		}
		if !Modules.Checker.Name(c, form.Value.(string)) {
			return
		}
		if Service.Checker.NameExist("biz", form.Value.(string), nil) {
			return
		}
	case "pic":
		if _, ok := form.Value.([]string); !ok {
			Modules.CallBack.Error(c, 101)
			return
		}
		if !Modules.Checker.Pic(c, form.Value.([]string)) {
			return
		}
		form.Value, _ = json.Marshal(form.Value)
	case "cat_address_id":
		if _, ok := form.Value.(uint); !ok {
			Modules.CallBack.Error(c, 101)
			return
		}
		if !Service.Checker.CatIdExist("address", form.Value.(uint)) {
			return
		}
	case "grade":
		if _, ok := form.Value.(uint); !ok {
			Modules.CallBack.Error(c, 101)
			return
		}
		if !Modules.Checker.Grade(c, form.Value.(uint)) {
			return
		}
	case "address":
		if _, ok := form.Value.(string); !ok {
			Modules.CallBack.Error(c, 101)
			return
		}
		if !Modules.Checker.Address(c, form.Value.(string)) {
			return
		}
	case "grade_weight":
		fallthrough
	case "max_price":
		fallthrough
	case "min_price":
		fallthrough
	case "husk":
		fallthrough
	case "share":
		fallthrough
	case "favorite":
		if _, ok := form.Value.(uint); !ok {
			Modules.CallBack.Error(c, 101)
			return
		}
	default:
		Modules.CallBack.Error(c, 114)
		return
	}
	if _, err := Service.Update(c, "biz", map[string]interface{}{
		form.Target: form.Value,
	}, map[string]interface{}{
		"id": form.Id,
	}); err != nil {
		return
	}
	Modules.CallBack.Default(c)
}

func (*biz) Delete(c *gin.Context) {
	var form struct {
		Id uint `form:"id" binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Service.Checker.BizExist(form.Id) {
		Modules.CallBack.Error(c, 116)
		return
	}
	//DEMO 删除菜单
	if _, err := Service.Delete(c, "biz", map[string]interface{}{
		"id": form.Id,
	}); err != nil {
		return
	}
	Modules.CallBack.Default(c)
}

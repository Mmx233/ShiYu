package Controllers

import (
	"Mmx/Modules"
	"Mmx/Service"
	"github.com/gin-gonic/gin"
	"strings"
)

type user struct{}

var User user

func (*user) Information(c *gin.Context) {
	username := c.Param("username")
	if !Service.Checker.AccountExist("user", username) {
		Modules.CallBack.Error(c, 111)
		return
	}
	type temp struct {
		Name      string `json:"name"`
		HeadImg   string `json:"head_img"`
		BigPlayer bool   `json:"big_player"`
		TestCount uint   `json:"test_count"`
		LikeCount uint16 `json:"like_count"`
		Favorites []uint `json:"favorites"`
		FanCount  string `json:"fan_count"`
		Fans      []uint `json:"fans"`
	}
	var t temp
	if Service.GetRow(c, "user", &t, map[string]interface{}{
		"username": username,
	}) != nil {
		return
	}
	Modules.Tool.ImgString(c, "user", &t.HeadImg) //处理头像url
	Modules.CallBack.Success(c, t)
}

func (*user) Renew(c *gin.Context) {
	username := c.Param("username")
	var form struct {
		UserName  string `form:"username" binding:"required,max=19"`
		PassWord  string `form:"password" binding:"required"`
		Name      string `form:"name" binding:"required"`
		BigPlayer bool   `form:"big_player"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Modules.Checker.Form(c, &form) {
		return
	}
	if username != form.UserName && Service.Checker.AccountExist("user", form.UserName) {
		Modules.CallBack.Error(c, 109)
		return
	}
	if !Service.Checker.Name(c, "user", username, form.Name) {
		return
	}
	salt := Modules.Tool.MakeSalt(form.PassWord)
	var insertMap = map[string]interface{}{
		"username": form.UserName,
		"password": Modules.Tool.EncodePassWord(form.PassWord, salt),
		"salt":     salt,
		"name":     form.Name,
	}
	if c.GetString("role") == "admin" { //admin可以修改是否为大玩家
		insertMap["big_player"] = form.BigPlayer
	}
	Modules.CallBack.Default(c)
}

func (*user) Change(c *gin.Context) {
	username := c.Param("username")
	var form struct {
		Target string      `form:"target" binding:"required"`
		Value  interface{} `form:"value"  binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	form.Target = strings.ToLower(form.Target)
	switch form.Target {
	case "username":
		if _, ok := form.Value.(string); !ok {
			Modules.CallBack.Error(c, 101)
			return
		}
		if !Modules.Checker.UserName(c, form.Value.(string)) {
			return
		}
		if Service.Checker.AccountExist("user", form.Value.(string)) {
			Modules.CallBack.Error(c, 109)
			return
		}
		if _, err := Service.Update(c, "user", map[string]interface{}{
			"username": form.Value,
		}, map[string]interface{}{
			"username": username,
		}); err != nil {
			return
		}
	case "password":
		if _, ok := form.Value.(string); !ok {
			Modules.CallBack.Error(c, 101)
			return
		}
		if !Modules.Checker.Password(c, form.Value.(string)) {
			return
		}
		salt := Modules.Tool.MakeSalt(form.Value.(string))
		if _, err := Service.Update(c, "user", map[string]interface{}{
			"password": Modules.Tool.EncodePassWord(form.Value.(string), salt),
			"salt":     salt,
		}, map[string]interface{}{
			"username": username,
		}); err != nil {
			return
		}
	case "name":
		if _, ok := form.Value.(string); !ok {
			Modules.CallBack.Error(c, 101)
			return
		}
		if !Modules.Checker.Name(c, form.Value.(string)) {
			return
		}
		if !Service.Checker.Name(c, "user", username, form.Value.(string)) {
			return
		}
		if _, err := Service.Update(c, "user", map[string]interface{}{
			"name": form.Value,
		}, map[string]interface{}{
			"username": username,
		}); err != nil {
			return
		}
	case "big_player":
		if _, ok := form.Value.(bool); !ok {
			Modules.CallBack.Error(c, 101)
			return
		}
		if c.GetString("role") == "admin" {
			if _, err := Service.Update(c, "user", map[string]interface{}{
				"big_player": form.Value.(bool),
			}, map[string]interface{}{
				"username": username,
			}); err != nil {
				return
			}
			break
		}
		fallthrough
	default:
		Modules.CallBack.Error(c, 114)
		return
	}
	Modules.CallBack.Default(c)
}

func (*user) Delete(c *gin.Context) {
	username := c.Param("username")
	if !Service.Checker.AccountExist("user", username) {
		Modules.CallBack.Error(c, 111)
		return
	}
	if _, err := Service.Delete(c, "user", map[string]interface{}{
		"username": username,
	}); err != nil {
		return
	}
	//删除头像图片
	//DEMO
	Modules.CallBack.Default(c)
}

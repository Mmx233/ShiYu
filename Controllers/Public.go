package Controllers

import (
	"Mmx/Modules"
	"Mmx/Service"
	"github.com/gin-gonic/gin"
)

type public struct{}

var Public public

func (*public) Login(c *gin.Context) {
	var form struct {
		Role     string `form:"role" binding:"required"`
		UserName string `form:"username" binding:"required,max=19"`
		PassWord string `form:"password" binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Modules.Checker.Form(c, &form) {
		return
	}
	if !Service.Checker.AccountExist(form.Role, form.UserName) {
		Modules.CallBack.Error(c, 111)
		return
	}
	if !Service.Checker.Password(c, form.Role, form.UserName, form.PassWord) {
		return
	}
	//登陆成功
	if token, err := Modules.Jwt.Encode(c, form.Role, form.UserName); err != nil {
		return
	} else {
		Modules.Cookie.SetCookie(c, "token", token)
	}
	Modules.CallBack.Default(c)
}

func (*public) Register(c *gin.Context) {
	var form struct {
		UserName string `form:"username" binding:"required,max=19"`
		PassWord string `form:"password" binding:"required"`
		Name     string `form:"name" binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Modules.Checker.Form(c, &form) {
		return
	}
	if Service.Checker.AccountExist("user", form.UserName) {
		Modules.CallBack.Error(c, 109)
		return
	}
	if Service.Checker.NameExist("user", form.Name) {
		Modules.CallBack.Error(c, 113)
		return
	}
	salt := Modules.Tool.MakeSalt(form.PassWord)
	if _, err := Service.Insert(c, "user", map[string]interface{}{
		"username":   form.UserName,
		"password":   Modules.Tool.EncodePassWord(form.PassWord, salt),
		"salt":       salt,
		"name": form.Name,
		"head_img":   "n",
		"big_player": 0,
		"test_count": 0,
		"like_count": 0,
		"fan_count":  0,
	}); err != nil {
		return
	}
	Modules.CallBack.Default(c)
}

package Controllers

import (
	"Mmx/Modules"
	"Mmx/Service"
	"github.com/gin-gonic/gin"
	"strings"
)

type admin struct{}

var Admin admin

func (*admin) NewAccount(c *gin.Context) {
	username := c.Param("username")
	var form struct {
		PassWord string `form:"password" binding:"required"`
		Name     string `form:"name" binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Modules.Checker.Form(c, &form) {
		return
	}
	if Service.Checker.AccountExist("admin", username) {
		Modules.CallBack.Error(c, 109)
		return
	}
	if Service.Checker.NameExist("admin", form.Name,nil) {
		Modules.CallBack.Error(c, 113)
		return
	}
	salt := Modules.Tool.MakeSalt(form.PassWord)
	if _, err := Service.Insert(c, "admin", map[string]interface{}{
		"username": username,
		"password": Modules.Tool.EncodePassWord(form.PassWord, salt),
		"name":     form.Name,
		"salt":     salt,
	}); err != nil {
		return
	}
	Modules.CallBack.Default(c)
}

func (*admin) Information(c *gin.Context) {
	username := c.Param("username")
	type adminInfo struct { //信息结构体
		Name string `json:"name"`
	}
	var admin adminInfo
	if Service.GetRow(c, "admin", &admin, map[string]interface{}{
		"username": username,
	}) != nil {
		return
	}
	Modules.CallBack.Success(c, admin)
}

func (*admin) Renew(c *gin.Context) {
	username := c.Param("username")
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
	if !Service.Checker.AccountExist("admin",username){
		Modules.CallBack.Error(c,111)
		return
	}
	if username != form.UserName && Service.Checker.AccountExist("admin", form.UserName) {
		Modules.CallBack.Error(c, 109)
		return
	}
	if !Service.Checker.Name(c, "admin", username, form.Name) {
		return
	}
	salt := Modules.Tool.MakeSalt(form.PassWord)
	if _, err := Service.Update(c, "admin", map[string]interface{}{
		"username": form.UserName,
		"password": Modules.Tool.EncodePassWord(form.PassWord, salt),
		"name":     form.Name,
		"salt":     salt,
	}, map[string]interface{}{
		"username": username,
	}); err != nil {
		return
	}
	Modules.CallBack.Default(c)
}

func (*admin) Change(c *gin.Context) {
	username := c.Param("username")
	var form struct {
		Target string `form:"target" binding:"required"`
		Value  string `form:"value"  binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Service.Checker.AccountExist("admin",username){
		Modules.CallBack.Error(c,111)
		return
	}
	form.Target = strings.ToLower(form.Target)
	switch form.Target {
	case "username":
		if !Modules.Checker.UserName(c, form.Value) {
			return
		}
		if Service.Checker.AccountExist("admin", form.Value) {
			Modules.CallBack.Error(c, 109)
			return
		}
		if _, err := Service.Update(c, "admin", map[string]interface{}{
			"username": form.Value,
		}, map[string]interface{}{
			"username": username,
		}); err != nil {
			return
		}
	case "password":
		if !Modules.Checker.Password(c, form.Value) {
			return
		}
		salt := Modules.Tool.MakeSalt(form.Value)
		if _, err := Service.Update(c, "admin", map[string]interface{}{
			"password": Modules.Tool.EncodePassWord(form.Value, salt),
			"salt":     salt,
		}, map[string]interface{}{
			"username": username,
		}); err != nil {
			return
		}
	case "name":
		if !Modules.Checker.Name(c, form.Value) {
			return
		}
		if !Service.Checker.Name(c, "admin", username, form.Value) {
			return
		}
		if _, err := Service.Update(c, "admin", map[string]interface{}{
			"name": form.Value,
		}, map[string]interface{}{
			"username": username,
		}); err != nil {
			return
		}
	default:
		Modules.CallBack.Error(c, 114)
		return
	}
	Modules.CallBack.Default(c)
}

func (*admin) Delete(c *gin.Context) {
	username := c.Param("username")
	if !Service.Checker.AccountExist("admin", username) {
		Modules.CallBack.Error(c, 111)
		return
	}
	if _, err := Service.Delete(c, "admin", map[string]interface{}{
		"username": username,
	}); err != nil {
		return
	}
	Modules.CallBack.Default(c)
}

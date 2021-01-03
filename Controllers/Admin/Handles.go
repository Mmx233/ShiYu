package Controllers

import (
	"Mmx/Modules"
	"Mmx/Service"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

type admin struct {}
var Admin admin

func (*admin)NewAccount(c *gin.Context){
	username:=c.Param("username")
	type Form struct {
		PassWord string `json:"password" form:"password" binding:"required"`
		Name string `json:"name" form:"name" binding:"required"`
	}
	var form Form
	if !Modules.Tool.BindForm(c, &form){
		return
	}
	if !Modules.Checker.Name(c,form.Name) {
		return
	}
	if !Modules.Checker.Password(c,form.PassWord){
		return
	}
	if Service.Checker.AccountExist("admin",username){
		Modules.CallBack.Error(c,109)
		return
	}
	salt:=Modules.Tool.MakeSalt(form.PassWord)
	if _,err:=Service.Insert(c,"admin",map[string]interface{}{
		"username":username,
		"password":Modules.Tool.EncodePassWord(form.PassWord,salt),
		"name":form.Name,
		"salt":salt,
	});err!=nil{
		return
	}
	Modules.CallBack.Default(c)
}

func (*admin)Information(c *gin.Context){
	username:=c.Param("username")
	type adminInfo struct {//信息结构体
		Name string `json:"name"`
	}
	var admin adminInfo
	if Service.GetRow(c,"admin",&admin,map[string]interface{}{
		"username":username,
	})!=nil{
		return
	}
	Modules.CallBack.Success(c,admin)
}

func (*admin)Renew(c * gin.Context){
	username:=c.Param("username")
	type renewForm struct {
		UserName string `json:"username" form:"username" binding:"required,max=19"`
		PassWord string `json:"password" form:"password" binding:"required"`
		Name string `json:"name" form:"name" binding:"required"`
	}
	var form renewForm
	if !Modules.Tool.BindForm(c,&form){
		return
	}
	if !Modules.Checker.Password(c,form.PassWord){
		return
	}
	if !Modules.Checker.Name(c,form.Name){
		return
	}
	if !Modules.Checker.UserName(c,form.UserName){
		return
	}
	if username!=form.UserName&&!Service.Checker.AccountExist("admin",form.UserName){
		Modules.CallBack.Error(c,111)
	}
	salt:=Modules.Tool.MakeSalt(form.PassWord)
	if _,err:=Service.Update(c,"admin",map[string]interface{}{
		"username":form.UserName,
		"password":Modules.Tool.EncodePassWord(form.PassWord,salt),
		"name":form.Name,
		"salt":salt,
	},map[string]interface{}{
		"username":username,
	});err!=nil{
		return
	}
	Modules.CallBack.Default(c)
}

func (*admin)Change(c *gin.Context){
	username:=c.Param("username")
	type changeForm struct {
		Target string `json:"target" form:"target" binding:"required"`
		Value string `json:"value" form:"value"  binding:"required"`
	}
	var form changeForm
	if !Modules.Tool.BindForm(c,&form){
		return
	}
	form.Target=strings.ToLower(form.Target)
	switch form.Target {
	case "username":
		if !Modules.Checker.UserName(c,form.Value){
			return
		}
		if Service.Checker.AccountExist("admin",form.Value){
			Modules.CallBack.Error(c,109)
			return
		}
		if _,err:=Service.Update(c,"admin",map[string]interface{}{
			"username":form.Value,
		}, map[string]interface{}{
			"username":username,
		});err!=nil{
			return
		}
	case "password":
		if !Modules.Checker.Password(c,form.Value){
			return
		}
		salt:=Modules.Tool.MakeSalt(form.Value)
		if _,err:=Service.Update(c,"admin", map[string]interface{}{
			"password":Modules.Tool.EncodePassWord(form.Value,salt),
			"salt":salt,
		}, map[string]interface{}{
			"username":username,
		});err!=nil{
			return
		}
	case "name":
		if !Modules.Checker.Name(c,form.Value){
			return
		}
		if _,err:=Service.Update(c,"admin", map[string]interface{}{
			"name":form.Value,
		}, map[string]interface{}{
			"username":username,
		});err!=nil{
			return
		}
	default:
		Modules.CallBack.ErrorWithErr(c,102,errors.New("目标属性不存在"))
		return
	}
	Modules.CallBack.Default(c)
}

func (*admin)Delete(c *gin.Context){
	username:=c.Param("username")
	if !Service.Checker.AccountExist("admin",username){
		Modules.CallBack.Error(c,111)
		return
	}
	if _,err:=Service.Delete(c,"admin", map[string]interface{}{
		"username":username,
	});err!=nil{
		return
	}
	Modules.CallBack.Default(c)
}
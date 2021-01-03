package Controllers

import (
	"Mmx/Modules"
	"Mmx/Service"
	"github.com/gin-gonic/gin"
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
	if Service.Check.AccountExist("admin",username){
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
	if !Service.Check.AccountExist("admin",form.UserName){
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

}
package Controllers

import (
	"Mmx/Modules"
	"Mmx/Service"
	"github.com/gin-gonic/gin"
)

type public struct {}
var Public public

func (*public)Login(c *gin.Context){
	type loginForm struct {
		Role string `json:"role" form:"role" binding:"required"`
		UserName string `json:"username" form:"username" binding:"required,max=19"`
		PassWord string `json:"password" form:"password" binding:"required"`
	}
	var form loginForm
	if !Modules.Tool.BindForm(c,&form){
		return
	}
	if !Modules.Checker.Form(c,&form){
		return
	}
	if !Service.Checker.AccountExist(form.Role,form.UserName){
		Modules.CallBack.Error(c,111)
		return
	}
	if !Service.Checker.Password(c,form.Role,form.UserName,form.PassWord){
		return
	}
	//登陆成功
	if token,err:=Modules.Jwt.Encode(c,form.Role,form.UserName);err!=nil{
		return
	}else{
		Modules.Cookie.SetCookie(c,"token",token)
	}
	Modules.CallBack.Default(c)
}

func (*public)Register(c *gin.Context){
	type registerForm struct {
		UserName string `json:"username" form:"username" binding:"required,max=19"`
		PassWord string `json:"password" form:"password" binding:"required"`
		Name string `json:"name" form:"name" binding:"required"`
	}
	var form registerForm
	if !Modules.Tool.BindForm(c,&form){
		return
	}
	if !Modules.Checker.Form(c,&form){
		return
	}

}
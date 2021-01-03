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
	if !Modules.Checker.Role(c,form.Role){
		return
	}
	if !Modules.Checker.UserName(c,form.UserName){
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

}

func (*public)Register(c *gin.Context){

}
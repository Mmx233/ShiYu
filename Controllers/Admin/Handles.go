package Controllers

import (
	"Mmx/Modules"
	"Mmx/Service"
	"github.com/gin-gonic/gin"
)

type admin struct {}
var Admin admin

func (*admin)NewAccount(c *gin.Context){
	type Form struct {
		PassWord string `json:"password" form:"password" binding:"required,len=32"`
		Name string `json:"name" form:"name" binding:"required"`
	}
	var form Form
	if err:=c.ShouldBind(&form);err!=nil{
		Modules.CallBack.ErrorWithErr(c,102,err)
		return
	}
	if !Modules.Checker.Name(c,form.Name) {//检查昵称
		return
	}
	if !Modules.Checker.Password(c,form.PassWord){//检查密码
		return
	}
	if Service.Check.AccountExist("admin",c.Param("username")){//用户名是否已被占用
		Modules.CallBack.Error(c,109)
		return
	}
	salt:=Modules.Tool.MakeSalt(form.PassWord)//生产盐
	if _,err:=Service.Insert(c,"admin",map[string]interface{}{
		"username":c.Param("username"),
		"password":Modules.Tool.EncodePassWord(form.PassWord,salt),
		"name":form.Name,
		"salt":salt,
	});err!=nil{
		return
	}
	Modules.CallBack.Default(c)
}

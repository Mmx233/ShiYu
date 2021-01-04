package Modules

import (
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

type checker struct {}
var Checker checker

func (*checker)Form(c *gin.Context,form interface{})bool{//利用反射自动检查表单，data必须传入结构体指针
	f:=reflect.ValueOf(form).Elem()
	g:=reflect.TypeOf(form).Elem()
	for i:=0;i<g.NumField();i++{
		switch strings.ToLower(g.Field(i).Name) {
		case "username":
			if !Checker.UserName(c,f.Field(i).String()){
				return false
			}
		case "name":
			if !Checker.Name(c,f.Field(i).String()){
				return false
			}
		case "password":
			if !Checker.Password(c,f.Field(i).String()){
				return false
			}
		case "role":
			if !Checker.Role(c,f.Field(i).String()){
				return false
			}
		}
	}
	return true
}

func (*checker)UserName(c *gin.Context,id string)bool{//检查用户名是否合法
	temp,err :=strconv.Atoi(id)
	if err!=nil  || len(id)>19 ||temp<1 {
		CallBack.Error(c,104)
		return false
	}
	return true
}

func (*checker)Name(c *gin.Context,content string)bool{//检查昵称是否合法
	if utf8.RuneCountInString(content)>8||len(content)==0{
		CallBack.Error(c, 103)
		return false
	}
	return true
}

func (*checker)Password(c *gin.Context,content string)bool{//检查密码是否合法
	if utf8.RuneCountInString(content)<9||utf8.RuneCountInString(content)>50{
		CallBack.Error(c,110)
		return false
	}
	return true
}

func (*checker)Role(c *gin.Context,role string)bool{
	if role!="admin"&&role!="user"{
		CallBack.ErrorWithErr(c,102,errors.New("role参数不合法"))
		return false
	}
	return true
}
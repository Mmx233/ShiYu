package Modules

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"unicode/utf8"
)

type checker struct {}
var Checker checker

func (*checker)UserName(c *gin.Context,id string)bool{
	temp,err :=strconv.Atoi(id)
	if err!=nil  || len(id)>19 ||temp<1 {
		CallBack.Error(c,104)
		return false
	}
	return true
}

func (*checker)Name(c *gin.Context,content string)bool{
	if utf8.RuneCountInString(content)>8||len(content)==0{
		CallBack.Error(c, 103)
		return false
	}
	return true
}

func (*checker)Password(c *gin.Context,content string)bool{
	if utf8.RuneCountInString(content)<9||utf8.RuneCountInString(content)>50{
		CallBack.Error(c,110)
		return false
	}
	return true
}
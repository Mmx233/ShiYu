package Modules

import (
	"Mmx/Service"
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

type checker struct{}

var Checker checker

func (*checker) Form(c *gin.Context, form interface{}) bool { //利用反射自动检查表单，data必须传入结构体指针
	f := reflect.ValueOf(form).Elem()
	g := reflect.TypeOf(form).Elem()
	for i := 0; i < g.NumField(); i++ {
		switch strings.ToLower(g.Field(i).Name) {
		case "username":
			if !Checker.UserName(c, f.Field(i).String()) {
				return false
			}
		case "name":
			if !Checker.Name(c, f.Field(i).String()) {
				return false
			}
		case "password":
			if !Checker.Password(c, f.Field(i).String()) {
				return false
			}
		case "role":
			if !Checker.Role(c, f.Field(i).String()) {
				return false
			}
		case "pic":
			if !Checker.Pic(c, f.Field(i).Interface().([]string)) {
				return false
			}
		case "address":
			if !Checker.Address(c, f.Field(i).String()) {
				return false
			}
		case "grade":
			if !Checker.Grade(c, f.Field(i).Interface().(uint)) {
				return false
			}
		case "cat":
			if !Checker.Cat(c, f.Field(i).String()) {
				return false
			}
		}
	}
	return true
}

func (*checker) Address(c *gin.Context, addr string) bool {
	if utf8.RuneCountInString(addr) > 225 {
		CallBack.Error(c, 118)
		return false
	}
	return true
}

func (*checker) UserName(c *gin.Context, id string) bool { //检查用户名是否合法
	temp, err := strconv.Atoi(id)
	if err != nil || len(id) > 19 || temp < 1 {
		CallBack.Error(c, 104)
		return false
	}
	return true
}

func (*checker) Name(c *gin.Context, content string) bool { //检查昵称是否合法
	if utf8.RuneCountInString(content) > 8 || len(content) == 0 {
		CallBack.Error(c, 103)
		return false
	}
	return true
}

func (*checker) Password(c *gin.Context, content string) bool { //检查密码是否合法
	if utf8.RuneCountInString(content) < 9 || utf8.RuneCountInString(content) > 50 {
		CallBack.Error(c, 110)
		return false
	}
	return true
}

func (*checker) Role(c *gin.Context, role string) bool {
	if role != "admin" && role != "user" {
		CallBack.ErrorWithErr(c, 102, errors.New("role参数不合法"))
		return false
	}
	return true
}

func (*checker) Pic(c *gin.Context, pics []string) bool {
	if pics == nil {
		CallBack.Error(c, 117)
		return false
	}
	for _, v := range pics {
		//不合规规则 DEMO
		/*if {
			CallBack.Error(c,117)
			return false
		}*/
	}
	return true
}

func (*checker) Grade(c *gin.Context, grade uint) bool {
	if grade > 10 {
		CallBack.Error(c, 121)
		return false
	}
	return true
}

func (*checker) Cat(c *gin.Context, cat string) bool {
	if cat != "address" && cat != "food" {
		CallBack.Error(c, 122)
		return false
	}
	return true
}
package Modules

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"reflect"
	"time"
	"unicode/utf8"
)

type tool struct {}
var Tool tool

func (*tool)EncodePassWord(p string,salt string)string{
	h:=sha256.New()
	h.Write([]byte(p+salt))
	r:=h.Sum(nil)
	return fmt.Sprintf("%x",r)
}

func randRune()string{
	rand.Seed(time.Now().UnixNano())
	letters:="abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	return string(letters[rand.Intn(len(letters))])
}

func (*tool)MakeSalt(PassWord string)string{
	Len:=utf8.RuneCountInString(PassWord)
	var salt string
	for i:=0;i<Len;i++{
		salt+=randRune()
	}
	return salt
}

func (*tool)BindForm(c *gin.Context,f interface{})bool{
	if err:=c.ShouldBind(f);err!=nil{
		CallBack.ErrorWithErr(c,102,err)
		return false
	}
	return true
}

func(*tool)ImgString(c *gin.Context,t string,d interface{}){
	switch t {
	case "biz":
		//商家图片
		rule:="https://xxx/%s"//商家图片规则 DEMO
		switch reflect.TypeOf(d).Elem().Kind(){
		case reflect.Slice:
			pics:=d.(*[]string)
			for i,v :=range *pics{
				(*pics)[i]=fmt.Sprintf(rule,v)
			}
		case reflect.String:
			pic:=d.(*string)
			*pic=fmt.Sprintf(rule,*pic)
		}
	case "user":
		//用户头像
		b:=d.(*string)
		if *b=="y"{
			u,ok:=c.Get("")
			if ok {
				*b = fmt.Sprintf("https://xx/%s", u) //头像规则 DEMO
				break
			}
		}
		*b="default"//默认头像
	}
}

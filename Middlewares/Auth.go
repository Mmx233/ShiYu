package Middlewares

import (
	"Mmx/Modules"
	"github.com/gin-gonic/gin"
	"net/http"
)

type auth struct {}
var Auth auth

func (* auth)Main(c *gin.Context){//鉴权中间件
	//防盗链
	if len(c.GetHeader("Referer"))<29||c.GetHeader("Referer")[8:28]!="hackweek.multmax.top"{
		Modules.CallBack.Error(c,302)
		return
	}

	sToken,err:=c.Cookie("token")
	if err!=nil||sToken==""{
		if c.FullPath()==`/api/login`||c.FullPath()==`/api/register`{
			c.Next()
			return
		}
		c.AsciiJSON(401,NewErrorCall("未登录",107))
		c.Abort()
		return
	}
	//解析jwt
	jToken, err := jwt.ParseWithClaims(sToken, &jwtDATA{}, func(token *jwt.Token) (interface{},error) {
		return jwtSEC, nil
	})
	claims, ok := jToken.Claims.(*jwtDATA)
	if err!=nil || !ok{
		c.AsciiJSON(500,NewErrorCall("鉴权token解析失败",106))
		c.Abort()
		return
	}
	if !jToken.Valid || claims.Ip!=c.ClientIP(){
		removeCookie(c,"token")
		c.AsciiJSON(401,NewErrorCall("登录过期",107))
		c.Abort()
		return
	}
	c.Set("username",claims.Username)
	if claims.Role!="admin" {
		switch{
		case "/api/user/" + claims.Username == c.Request.RequestURI:
			fallthrough
		case c.Request.Method=="GET"&&c.Request.RequestURI[0:9]=="/api/biz/"://允许biz GET
			fallthrough
		case c.FullPath()=="/api/register"||c.FullPath()=="/api/login"://登录注册页
			fallthrough
		case c.FullPath()=="/api/biz/menu/"+claims.Username://收藏
			fallthrough
		case c.FullPath()=="/api/upload/3"://头像上传
			c.Next()
			return
		}
		c.AsciiJSON(403, NewErrorCall("没有权限", 108))
		c.Abort()
		return
	}
	c.Next()
}

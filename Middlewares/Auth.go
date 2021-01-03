package Middlewares

import (
	"Mmx/Modules"
	"github.com/gin-gonic/gin"
)

type auth struct {}
var Auth auth

func (* auth)Main(c *gin.Context){//鉴权中间件
	sToken,err:=c.Cookie("token")
	if err!=nil||sToken==""{//未登录
		if c.FullPath()==`/api/login`||c.FullPath()==`/api/register`{//特殊入口免鉴权
			c.Next()
			return
		}
		Modules.CallBack.Error(c,107)
		return
	}

	//jwt
	var claims *Modules.JwtDATA
	{
		var err error
		if claims,err=Modules.Jwt.Decode(c,sToken);err!=nil{
			Modules.Cookie.RemoveCookie(c,"token")
			return
		}
	}

	//传递登录者username
	c.Set("username",claims.Username)

	/*if claims.Role!="admin" {//user权限 DEMO
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
		Modules.CallBack.Error(c,108)
		return
	}*/
	c.Next()
}

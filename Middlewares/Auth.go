package Middlewares

import (
	"Mmx/Modules"
	"github.com/gin-gonic/gin"
	"strings"
)

type auth struct{}

var Auth auth

func (*auth) Main(c *gin.Context) { //鉴权中间件
	sToken:= c.GetHeader("authorization")
	if c.FullPath() == `/api/v3/login` || c.FullPath() == `/api/v3/register` { //特殊入口免鉴权
		c.Next()
		return
	}
	if sToken == "" { //未登录
		Modules.CallBack.Error(c, 107)
		return
	}

	//jwt
	var claims *Modules.JwtDATA
	{
		var err error
		if claims, err = Modules.Jwt.Decode(c, sToken); err != nil {
			return
		}
	}

	//传递登录者username和role
	c.Set("username", claims.Username)
	c.Set("role", claims.Role)

	//user权限
	if claims.Role != "admin" {
		switch {
		case "/api/user/"+claims.Username == c.Request.RequestURI:
			fallthrough
		case strings.HasPrefix(c.FullPath(), "/api/biz/menu/fav/"): //用户收藏
			fallthrough
		case c.Request.Method == "GET" && strings.HasPrefix(c.FullPath(), "/api/biz/"): //允许biz GET
			fallthrough
		case strings.HasPrefix(c.FullPath(), "/api/v3/"): //公共api
			c.Next()
			return
		}
		Modules.CallBack.Error(c, 108)
		return
	}
	c.Next()
}

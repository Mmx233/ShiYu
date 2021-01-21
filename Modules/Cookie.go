package Modules

import "github.com/gin-gonic/gin"

type cookie struct{}

var Cookie cookie

func (*cookie) SetCookie(c *gin.Context, name string, value string) {
	c.SetCookie(name, value, 2590000, "/", "", true, true)
}

func (*cookie) RemoveCookie(c *gin.Context, name string) {
	c.SetCookie(name, "", -1, "/", "", true, true)
}

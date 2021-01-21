package Middlewares

import (
	"Mmx/Modules"
	"github.com/gin-gonic/gin"
)

type check struct{}

var Check check

func (*check) UserName(c *gin.Context) {
	if !Modules.Checker.UserName(c, c.Param("username")) {
		return
	}
	c.Next()
}

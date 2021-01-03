package Service

import (
	"Mmx/Modules"
	"github.com/gin-gonic/gin"
)

type check struct {}
var Checker check

func (*check)AccountExist(role string,username string)bool{
	var tempB bool
	DB.QueryRow("SELECT 1 FROM "+role+" WHERE username=?",username).Scan(&tempB)
	return tempB
}

func (*check)Password(c *gin.Context,Role string,UserName string,PassWord string)bool{
	type p struct {
		PassWord string
		Salt string
	}
	var P p
	if GetRow(c,Role,&P, map[string]interface{}{
		"username":UserName,
	})!=nil{
		return false
	}
	if Modules.Tool.EncodePassWord(PassWord,P.Salt)!=P.PassWord{
		Modules.CallBack.Error(c,100)
		return false
	}
	return true
}
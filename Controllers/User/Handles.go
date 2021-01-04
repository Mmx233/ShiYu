package Controllers

import (
	"Mmx/Modules"
	"Mmx/Service"
	"github.com/gin-gonic/gin"
)

type user struct {}
var User user

func (*user)Information(c *gin.Context){
	username:=c.Param("username")
	if !Service.Checker.AccountExist("user",username){
		Modules.CallBack.Error(c,111)
		return
	}
	type temp struct {
		Name string `json:"name"`
		HeadImgId uint `json:"head_img_id"`
		BigPlayer bool `json:"big_player"`
		TestCount uint `json:"test_count"`
		LikeCount uint16 `json:"like_count"`
		Favorites []uint `json:"favorites"`
		FanCount string `json:"fan_count"`
		Fans []uint `json:"fans"`
	}
	var t temp
	if Service.GetRow(c,"user",&t, map[string]interface{}{
		"username":username,
	})!=nil{
		return
	}
	Modules.CallBack.Success(c,t)
}
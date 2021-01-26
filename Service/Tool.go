package Service

import "github.com/gin-gonic/gin"

type tool struct {}
var Tool tool

func (*tool)GetFavorite(c *gin.Context,table string,id uint)(uint,error){
	//获取收藏数
	var d struct{
		Favorite uint `json:"favorite"`
	}
	err := GetRow(c, table, &d, map[string]interface{}{
		"id": id,
	})
	return d.Favorite,err
}

func (*tool)FavoritePlus(c *gin.Context,table string,id uint,plus int)error{
	num,err:=Tool.GetFavorite(c,table,id)
	if err!=nil{
		return err
	}
	if _,err:=Update(c,table, map[string]interface{}{
		"favorite":uint(int32(num)+int32(plus)),
	}, map[string]interface{}{
		"id":id,
	});err!=nil{
		return err
	}
	return nil
}

func (*tool)BizFavoritePlus(c *gin.Context,id uint,plus int)error{
	return Tool.FavoritePlus(c,"biz",id,plus)
}

func (*tool)MenuFavoritePlus(c *gin.Context,id uint,plus int)error{
	var d struct{
		BizId uint `json:"biz_id"`
	}
	if err:=GetRow(c,"biz_menu",&d, map[string]interface{}{
		"id":id,
	});err!=nil{
		return err
	}
	if err:=Tool.BizFavoritePlus(c,d.BizId,plus);err!=nil{
		return err
	}
	return Tool.FavoritePlus(c,"biz_menu",id,plus)
}
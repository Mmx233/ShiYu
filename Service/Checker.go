package Service

import (
	"Mmx/Modules"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type check struct{}

var Checker check

func (*check) AccountExist(role string, username string) bool {
	var temp bool
	DB.QueryRow("SELECT 1 FROM "+role+" WHERE username=?", username).Scan(&temp)
	return temp
}

func (*check) NameExist(table string, name string, where map[string]interface{}) bool {
	var v []string
	var values []interface{}
	values=append(values,name)
	if where!=nil&&len(where)!=0{
		v,values=Modules.Tool.MakeWhere(where)
		v[0]=","+v[0]
	}
	var temp bool
	DB.QueryRow(fmt.Sprintf("SELECT 1 FROM "+table+" WHERE name=?%v",strings.Join(v,",")), values...).Scan(&temp)
	return temp
}

func (*check) BizExist(id uint) bool {
	var temp bool
	DB.QueryRow("SELECT 1 FROM biz WHERE id=?", id).Scan(&temp)
	return temp
}

func (*check) MenuExist(id uint) bool {
	var temp bool
	DB.QueryRow("SELECT 1 FROM biz_menu WHERE id=?", id).Scan(&temp)
	return temp
}

func (*check) Password(c *gin.Context, Role string, UserName string, PassWord string) bool { //登陆时的核对密码是否正确
	type p struct {
		PassWord string
		Salt     string
	}
	var P p
	if GetRow(c, Role, &P, map[string]interface{}{
		"username": UserName,
	}) != nil {
		return false
	}
	if Modules.Tool.EncodePassWord(PassWord, P.Salt) != P.PassWord {
		Modules.CallBack.Error(c, 100)
		return false
	}
	return true
}

func (*check) CatIdExist(cat string, id uint) bool {
	var temp bool
	DB.QueryRow("SELECT 1 FROM cat_"+cat+" WHERE id=?", id).Scan(&temp)
	return temp
}

func (*check) Name(c *gin.Context, role string, username string, name string) bool { //检查昵称是否被占用
	if Checker.NameExist(role, name,nil) {
		type temp struct {
			UserName string
			Name     string
		}
		var t temp
		if GetRow(c, role, &t, map[string]interface{}{
			"name": name,
		}) != nil {
			return false
		}
		if t.UserName != username {
			Modules.CallBack.Error(c, 113)
			return false
		}
	}
	return true
}

func (*check) IsFav(c *gin.Context, id uint) (bool, []uint) {
	//特殊函数，接收了c但是不写报错
	var u struct {
		Fav []uint `json:"favorites"`
	}
	if GetRow(nil, "user", &u, map[string]interface{}{
		"username": c.GetString("username"),
	}) != nil {
		return false, nil
	}
	return Modules.Tool.Find(u.Fav, id), u.Fav
}

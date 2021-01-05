package Modules

import "github.com/gin-gonic/gin"

type callBack struct {}
var CallBack callBack

//信息结构体
var errCode =map[int]int{//错误代码对应Http状态码
	100:401,
	101:400,
	102:500,
	103:400,
	104:400,
	105:401,
	106:500,
	107:401,
	108:403,
	109:409,
	110:400,
	111:404,
	112:500,
	113:409,
	114:404,
	301:429,
	302:403,
	303:429,
}

var errMsg =map[int]string{//错误码对应错误信息
	100:"密码错误",
	101:"参数错误",
	102:"数据库操作失败",
	103:"昵称不合法",
	104:"目标用户名不合法",
	105:"登录过期",
	106:"鉴权token解析失败",
	107:"未登录",
	108:"没有权限",
	109:"用户名已被占用",
	110:"新密码不合法",
	111:"目标账号不存在",
	112:"鉴权token生产失败",
	113:"昵称已被占用",
	114:"目标属性不存在",
	301:"您的访问过快，请稍后",
	302:"拒绝访问",
	303:"您的访问过快，请30分钟后再试",
}

type message struct {//callback结构
	Status string `json:"status"`
	Code int `json:"code"`
	Data interface{} `json:"data"`
	Msg string `json:"msg"`
}

func (* callBack)Success(c *gin.Context,data interface{}){
	c.AsciiJSON(200,message{
		Status: "success",
		Code:200,
		Data: data,
		Msg:"",
	})
}

func (* callBack)Error(c *gin.Context,code int){
	c.AsciiJSON(errCode[code],message{
		Status: "Error",
		Code: code,
		Data: map[int]string{},
		Msg: errMsg[code],
	})
	c.Abort()
}

func (* callBack)ErrorWithErr(c *gin.Context,code int,err error){
	c.AsciiJSON(errCode[code],message{
		Status: "Error",
		Code: code,
		Data: map[int]string{},
		Msg: err.Error(),
	})
}

func (*callBack)Default(c *gin.Context){
	c.AsciiJSON(200,message{
		Status: "Success",
		Code: 200,
		Data: map[int]string{},
		Msg: "",
	})
}
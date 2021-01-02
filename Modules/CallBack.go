package Modules

import "github.com/gin-gonic/gin"

type callBack struct {}
var CallBack callBack

//信息结构体
var errCode =map[int]int{//错误代码对应Http状态码
	100:401,
	301:429,
	302:403,
	303:429,
}

var errMsg =map[int]string{//错误码对应错误信息
	100:"密码错误",
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

func (* callBack)Success(c *gin.Context,data interface{},msg string){
	c.AsciiJSON(200,message{
		Status: "success",
		Code:200,
		Data: data,
		Msg:msg,
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
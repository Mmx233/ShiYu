package Middlewares

import (
	"Mmx/Modules"
	"github.com/gin-gonic/gin"
	"time"
)

type secure struct {}
var Sec secure

var (
	ipLogger=make(map[string]int)
	ipChan  = make(chan string)
)

func (* secure)InitIpLogger(){
	for{
		ip:=<-ipChan
		if ipLogger[ip]>=0 { //配合封禁
			ipLogger[ip]++
			go func(){//仅记录60秒内的访问
				time.Sleep(time.Minute)
				ipLogger[ip]--
				if ipLogger[ip]==0 {
					delete(ipLogger, ip)
				}
			}()
		}
	}
}

func (* secure)Main(c *gin.Context){
	//防扫描
	ipChan <-c.ClientIP()
	if ipLogger[c.ClientIP()]>60 || ipLogger[c.ClientIP()]<0 { //一分钟内最多60次访问，限制访问频次
		Modules.CallBack.Error(c,301)
		return
	}else if ipLogger[c.ClientIP()]>=240{ //每分钟超240次封禁IP
		ipLogger[c.ClientIP()]=-1 //使被拦截
		go func(ip string){//截除拦截
			time.Sleep(time.Hour/2)//半小时后截除
			delete(ipLogger,ip)
		}(c.ClientIP())
		Modules.CallBack.Error(c,303)
		return
	}
	c.Next()
}

package Middlewares

import (
	"Mmx/Modules"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type secure struct{}

var Sec secure

type decrease struct {
	Ip string
	Time int64
}

var (
	ipLogger = make(map[string]int)
	ipChan   = make(chan string)
	decreaseChan = make(chan decrease)
)

func (*secure) InitIpLogger() {
	go func () {//ip记录
		for {
			ip := <-ipChan
			if ipLogger[ip] >= 0 { //配合封禁
				ipLogger[ip]++
				decreaseChan <- decrease{
					ip,
					time.Now().Unix()+60,//60s后消除
				}
			}
		}
	}()

	go func() {//ip计数消除栈
		var d []decrease
		var HaveWorker bool
		var worker=func (){//消除执行
			HaveWorker=true
			for {
				if t := d[0].Time - time.Now().Unix(); t > 0 {
					time.Sleep(time.Duration(t) * time.Second)
				}
				ipLogger[d[0].Ip]--
				if (ipLogger[d[0].Ip]) == 0 {
					delete(ipLogger, d[0].Ip)
				}
				if len(d) > 1 {
					d = d[1:]
				} else {
					d = make([]decrease, 0)
					HaveWorker=false
					break
				}
			}
		}
		go func() {//将数据接收入栈
			for{
				data := <- decreaseChan
				d=append(d,data)
				if !HaveWorker{
					go worker()
				}
			}
		}()
	}()
}

func (*secure) Main(c *gin.Context) {
	//防盗链
	if c.GetHeader("Referer") != "" && strings.HasPrefix(c.GetHeader("Referer"), "https://shiyu.icu") {
		Modules.CallBack.Error(c, 302)
		return
	}
	ipChan <- c.ClientIP()
	if ipLogger[c.ClientIP()] > 60 || ipLogger[c.ClientIP()] < 0 { //一分钟内最多60次访问，限制访问频次
		Modules.CallBack.Error(c, 301)
		return
	} else if ipLogger[c.ClientIP()] >= 240 { //每分钟超240次封禁IP
		ipLogger[c.ClientIP()] = -1 //使被拦截
		go func(ip string) { //截除拦截
			time.Sleep(time.Hour / 2) //半小时后解除
			delete(ipLogger, ip)
		}(c.ClientIP())
		Modules.CallBack.Error(c, 303)
		return
	}
	c.Next()
}

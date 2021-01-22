package Modules

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"reflect"
	"sort"
	"time"
	"unicode/utf8"
)

type tool struct{}

var Tool tool

func (*tool) EncodePassWord(p string, salt string) string {
	h := sha256.New()
	h.Write([]byte(p + salt))
	r := h.Sum(nil)
	return fmt.Sprintf("%x", r)
}

func randRune() string {
	rand.Seed(time.Now().UnixNano())
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	return string(letters[rand.Intn(len(letters))])
}

func (*tool) MakeSalt(PassWord string) string {
	Len := 50-utf8.RuneCountInString(PassWord)
	var salt string
	for i := 0; i < Len; i++ {
		salt += randRune()
	}
	return salt
}

func (*tool) BindForm(c *gin.Context, f interface{}) bool {
	if err := c.ShouldBind(f); err != nil {
		CallBack.ErrorWithErr(c, 101, err)
		return false
	}
	return true
}

func (*tool) ImgString(c *gin.Context, t string, d interface{}) {
	switch t {
	case "biz":
		//商家图片
		rule := "https://xxx/%s" //商家图片规则 DEMO
		switch reflect.TypeOf(d).Elem().Kind() {
		case reflect.Slice:
			pics := d.(*[]string)
			for i, v := range *pics {
				(*pics)[i] = fmt.Sprintf(rule, v)
			}
		case reflect.String:
			pic := d.(*string)
			*pic = fmt.Sprintf(rule, *pic)
		}
	case "user":
		//用户头像
		b := d.(*string)
		if *b == "y" {
			u, ok := c.Get("")
			if ok {
				*b = fmt.Sprintf("https://xx/%s", u) //头像规则 DEMO
				break
			}
		}
		*b = "default" //默认头像
	}
}

func (*tool) Find(a []uint, b uint) bool { //find b in a
	var c []int
	d := int(b)
	for _, v := range a {
		c = append(c, int(v))
	}
	sort.Ints(c)
	if len(c) == 0 {
		return false
	}
	if c[0] == d || c[len(c)-1] == d {
		return true
	}
	var i = 0
	var ii = len(c) - 1
	for {
		iii := (ii - i) / 2
		if iii == 0 {
			break
		}
		if c[iii+i] == d {
			return true
		} else if a[iii+i] > b {
			ii = ii - iii
		} else {
			i += iii
		}
	}
	return false
}

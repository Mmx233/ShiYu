package Modules

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
	"unicode/utf8"
)

type tool struct {}
var Tool tool

func (*tool)EncodePassWord(p string,salt string)string{
	h:=sha256.New()
	h.Write([]byte(p+salt))
	r:=h.Sum(nil)
	return fmt.Sprintf("%x",r)
}

func randRune()string{
	rand.Seed(time.Now().UnixNano())
	letters:="abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	return string(letters[rand.Intn(len(letters))])
}

func (*tool)MakeSalt(PassWord string)string{
	Len:=utf8.RuneCountInString(PassWord)
	var salt string
	for i:=0;i<Len;i++{
		salt+=randRune()
	}
	return salt
}

package Modules

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

type JwtDATA struct {
	Role     string
	Username string
	jwt.StandardClaims
}

type j struct {
	sec []byte
}

func initJwt() j {
	return j{
		sec: []byte("我叼，密钥应该复杂一点，对吧对吧"),
	}
}

var Jwt = initJwt()

func (j *j) Decode(c *gin.Context, sToken string) (*JwtDATA, error) {
	jToken, err := jwt.ParseWithClaims(sToken, &JwtDATA{}, func(token *jwt.Token) (interface{}, error) {
		return j.sec, nil
	})
	claims, ok := jToken.Claims.(*JwtDATA)
	if err != nil || !ok {
		CallBack.Error(c, 106)
		return claims, err
	}
	if !jToken.Valid {
		CallBack.Error(c, 105)
		return claims, errors.New("登录过期")
	}
	return claims, nil
}

func (j *j) Encode(c *gin.Context, role string, username string) (string, error) {
	JToken := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtDATA{
		Role:     role,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(), // 过期时间
			Issuer:    "Mmx",                                      // 签发人
		},
	})
	token, err := JToken.SignedString(j.sec)
	if err != nil {
		CallBack.Error(c, 112)
	}
	return token, nil
}

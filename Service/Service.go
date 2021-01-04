package Service

import (
	"Mmx/Modules"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"reflect"
	"strings"
)

var DB *sql.DB

func er(c *gin.Context,e error){//向响应体写入报错json
	if e!=nil{
		Modules.CallBack.Error(c,102)
	}
}

func InitDatabase(){//连接数据库
	var err error
	if DB, err = sql.Open("mysql", "HackWeek:EfhGjswCa73efjeX@tcp(127.0.0.1:3306)/HackWeek");err!=nil||DB.Ping()!=nil{
		fmt.Println(DB.Ping())
		fmt.Println(err)
		os.Exit(3)
	}
}

func Exec(sql string,args ...interface{})(int64,error){
	if r,e:=DB.Exec(sql,args...);e!=nil{
		return 0, e
	}else{
		if id,er:=r.LastInsertId();er!=nil{
			return 0,er
		}else{
			return id,nil
		}
	}
}

func Insert(c *gin.Context,table string,value map[string]interface{})(int64,error){
	var keys []string
	var p []string
	var values []interface{}
	for k,v:=range value{
		keys=append(keys,k)
		p=append(p,"?")
		values=append(values,v)
	}
	SQL:=fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",table,strings.Join(keys,","),strings.Join(p,","))
	id,err:= Exec(SQL,values...)
	er(c,err)
	return id,err
}

func Update(c *gin.Context,table string,value map[string]interface{},where map[string]interface{})(int64,error){
	var keys []string
	var wh []string
	var values []interface{}
	for k,v:=range value{
		keys=append(keys,k+"=?")
		values=append(values,v)
	}
	for w,v:=range where{
		wh=append(wh,w)
		values=append(values,v)
	}
	SQL:=fmt.Sprintf("UPDATE %s SET %s WHERE %s",table,strings.Join(keys,","),strings.Join(wh,","))
	id,err:= Exec(SQL,values...)
	er(c,err)
	return id,err
}

func Delete(c *gin.Context,table string,where map[string]interface{})(int64,error){
	var wh []string
	var values []interface{}
	for w,v:=range where{
		wh=append(wh,w)
		values=append(values,v)
	}
	SQL:=fmt.Sprintf("DELETE FROM %s WHERE %s",table,strings.Join(wh,","))
	id,err:= Exec(SQL,values...)
	er(c,err)
	return id,err
}

func GetRow(c *gin.Context,table string,data interface{},where map[string]interface{})error{//传入的data为结构体的指针，否则无法工作
	var points []interface{}
	var keys []string
	var wh []string
	var values []interface{}
	{//获取data结构体信息
		f:=reflect.ValueOf(data).Elem()
		for i:=0;i<f.NumField();i++{//取址
			points=append(points, f.Field(i).Addr())
		}
		g:=reflect.TypeOf(data).Elem()
		for i:=0;i<g.NumField();i++{//获取对应字段名，对应数据库中键名
			keys=append(keys,strings.ToLower(g.Field(i).Name))
		}
	}
	for k,v :=range where{
		wh=append(wh,k+"=?")
		values=append(values,v)
	}
	SQL:=fmt.Sprintf("SELECT %s FROM %s WHRER %s",strings.Join(keys,","),table,strings.Join(wh,","))
	row:=DB.QueryRow(SQL,values...)
	if row.Err()!=nil{
		er(c,row.Err())
		return row.Err()
	}
	if err:=row.Scan(points...);err!=nil{
		er(c,err)
		return err
	}
	return nil
}
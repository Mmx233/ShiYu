package Service

import (
	"Mmx/Modules"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"reflect"
	"strings"
)

var DB *sql.DB

func er(c *gin.Context, e error) { //向响应体写入报错json
	if e != nil {
		Modules.CallBack.Error(c, 102)
	}
}

func InitDatabase() { //连接数据库
	var err error
	if DB, err = sql.Open("mysql", "HackWeek:EfhGjswCa73efjeX@tcp(127.0.0.1:3306)/HackWeek"); err != nil || DB.Ping() != nil {
		fmt.Println(DB.Ping())
		fmt.Println(err)
		os.Exit(3)
	}
}

func Exec(sql string, args ...interface{}) (int64, error) {
	if r, e := DB.Exec(sql, args...); e != nil {
		return 0, e
	} else {
		if id, er := r.LastInsertId(); er != nil {
			return 0, er
		} else {
			return id, nil
		}
	}
}

func Insert(c *gin.Context, table string, value map[string]interface{}) (int64, error) {
	var keys []string
	var p []string
	var values []interface{}
	for k, v := range value {
		keys = append(keys, k)
		p = append(p, "?")
		values = append(values, v)
	}
	SQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(keys, ","), strings.Join(p, ","))
	id, err := Exec(SQL, values...)
	er(c, err)
	return id, err
}

func Update(c *gin.Context, table string, value map[string]interface{}, where map[string]interface{}) (int64, error) {
	var keys []string
	var wh []string
	var values []interface{}
	for k, v := range value {
		keys = append(keys, k+"=?")
		values = append(values, v)
	}
	for w, v := range where {
		wh = append(wh, w)
		values = append(values, v)
	}
	SQL := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(keys, ","), strings.Join(wh, ","))
	id, err := Exec(SQL, values...)
	er(c, err)
	return id, err
}

func Delete(c *gin.Context, table string, where map[string]interface{}) (int64, error) {
	var wh []string
	var values []interface{}
	for w, v := range where {
		wh = append(wh, w)
		values = append(values, v)
	}
	SQL := fmt.Sprintf("DELETE FROM %s WHERE %s", table, strings.Join(wh, ","))
	id, err := Exec(SQL, values...)
	er(c, err)
	return id, err
}

func GetRow(c *gin.Context, table string, data interface{}, where map[string]interface{}) error { //传入的data为结构体的指针，否则无法工作
	var points []interface{}
	var keys []string
	var wh []string
	var values []interface{}

	//临时数组json 原字符串存储
	var temp = make(map[int]*string)

	//获取data结构体信息
	f := reflect.ValueOf(data).Elem()
	for i := 0; i < f.NumField(); i++ { //取址
		if f.Field(i).Kind() == reflect.Slice {
			//对应数据库的数组存储
			//目前数据库所有数组都为[]uint
			//对数组的还原在本函数底部
			var s string
			points = append(points, &s)
			temp[i] = &s
			continue
		}
		points = append(points, f.Field(i).Addr())
	}
	g := reflect.TypeOf(data).Elem()
	for i := 0; i < g.NumField(); i++ { //获取对应字段名，对应数据库中键名
		if g.Field(i).Tag.Get("json") != "" {
			keys = append(keys, strings.ToLower(g.Field(i).Tag.Get("json"))) //有json tag则使用
		} else {
			keys = append(keys, strings.ToLower(g.Field(i).Name))
		}
	}

	//构造查询语句
	for k, v := range where {
		wh = append(wh, k+"=?")
		values = append(values, v)
	}
	SQL := fmt.Sprintf("SELECT %s FROM %s WHRER %s", strings.Join(keys, ","), table, strings.Join(wh, ","))
	row := DB.QueryRow(SQL, values...)
	if row.Err() != nil {
		er(c, row.Err())
		return row.Err()
	}
	if err := row.Scan(points...); err != nil {
		er(c, err)
		return err
	}

	//还原数组
	for k, v := range temp {
		var tempS []uint64
		if json.Unmarshal([]byte(*v), &tempS) != nil {
			err := errors.New("未知错误-解码失败")
			Modules.CallBack.ErrorWithErr(c, 102, err)
			return err
		}
		for i, vv := range tempS {
			f.Field(k).Index(i).SetUint(vv) //默认全是[]uint
		}
	}

	return nil
}

func Get(c *gin.Context, table string, data interface{}, where map[string]interface{}, page int) error { //传入data为结构体的切片的指针，必须make，按容量扫描
	var keys []string
	var wh []string
	var values []interface{}

	var limit int

	{ //取字段名
		limit = reflect.ValueOf(data).Elem().Len()
		j := reflect.TypeOf(reflect.ValueOf(data).Elem().Index(0).Interface())
		for i := 0; i < j.NumField(); i++ {
			if j.Field(i).Tag.Get("json") != "" {
				keys = append(keys, j.Field(i).Tag.Get("json"))
			} else {
				keys = append(keys, strings.ToLower(j.Field(i).Name))
			}
		}
	}

	for k, v := range where {
		wh = append(wh, k+"=?")
		values = append(values, v)
	}
	page = (page - 1) * limit
	SQL := fmt.Sprintf("SELECT %s FROM %s WHERE %s limit %d,%d", strings.Join(keys, ","), table, strings.Join(wh, ","), page, page+limit)
	rows, err := DB.Query(SQL, values...)
	defer rows.Close()
	if err != nil {
		er(c, err)
		return err
	}

	//数组临时存储
	var a1 []*string
	var a2 []interface{}

	i := 0
	for rows.Next() {
		//取指针
		var points []interface{}
		o := reflect.ValueOf(reflect.ValueOf(data).Elem().Index(i).Interface())
		for ii := 0; ii < o.NumField(); ii++ {
			if o.Field(ii).Kind() == reflect.Slice {
				var s string
				points = append(points, &s)
				a1 = append(a1, &s)
				a2 = append(a2, o.Field(ii).Addr())
				continue
			}
			points = append(points, o.Field(ii).Addr())
		}

		//导入数据
		if err := rows.Scan(points...); err != nil {
			er(c, err)
			return err
		}

		i++
	}
	//应对数据不够一页及超出
	if i == 0 {
		Modules.CallBack.Error(c, 115)
		return errors.New("没有更多了")
	} else if i != limit {
		reflect.ValueOf(data).Elem().SetLen(i)
	}

	//数组还原
	for ii, v := range a1 {
		var temp []uint64
		if json.Unmarshal([]byte(*v), &temp) != nil {
			err := errors.New("未知错误-解码失败")
			Modules.CallBack.ErrorWithErr(c, 102, err)
			return err
		}
		for iii, vv := range temp {
			reflect.ValueOf(a2[ii]).Elem().Index(iii).SetUint(vv) //默认全是[]uint
		}
	}

	return nil
}

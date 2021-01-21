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
		//DEMO 应对数组
		/*switch reflect.TypeOf(v) {

		}*/
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

func getKeys(v reflect.Type) []string {
	var keys []string
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Tag.Get("skip") == "true" {
			continue
		}
		if v.Field(i).Tag.Get("json") != "" {
			keys = append(keys, v.Field(i).Tag.Get("json"))
			continue
		}
		keys = append(keys, strings.ToLower(v.Name()))
	}
	return keys
}

func getKeysAndPointers(d interface{}) ([]string, []interface{}, []*string, []interface{}) { //获取结构体字段名的字符串切片及其地址
	//返回的值依次为 字段切片 地址指针切片 临时数组字符串指针切片 数组真实指针切片
	var keys []string
	var pointers []interface{}
	var a1 []*string
	var a2 []interface{}
	v := reflect.TypeOf(d).Elem()
	keys = getKeys(v) //获取字段名
	f := reflect.ValueOf(d).Elem()
	for i := 0; i < f.NumField(); i++ { //获取其他
		if v.Field(i).Tag.Get("skip") == "true" {
			continue
		}
		if f.Field(i).Kind() == reflect.Slice {
			var s string
			pointers = append(pointers, &s)
			a1 = append(a1, &s)
			a2 = append(a2, f.Field(i).Addr())
		}
		pointers = append(pointers, f.Field(i).Addr())
	}
	return keys, pointers, a1, a2
}

func turnSliceBack(c *gin.Context, a1 []*string, a2 []interface{}) error {
	for ii, v := range a1 {
		var temp interface{}
		if strings.Contains(reflect.TypeOf(a2[ii]).Elem().String(), "\"") {
			temp = make([]string, 0)
		} else {
			temp = make([]uint64, 0)
		}
		if json.Unmarshal([]byte(*v), &temp) != nil {
			err := errors.New("未知错误-解码失败")
			Modules.CallBack.ErrorWithErr(c, 102, err)
			return err
		}
		e := make([]reflect.Value, 0)
		q := reflect.ValueOf(temp)
		for iii := 0; iii < q.NumField(); iii++ {
			e = append(e, q.Field(iii))
		}
		reflect.ValueOf(a2[ii]).Elem().Set(reflect.Append(reflect.ValueOf(a2[ii]), e...))
	}
	return nil
}

func GetRow(c *gin.Context, table string, data interface{}, where map[string]interface{}) error { //传入的data为结构体的指针，否则无法工作
	var wh []string
	var values []interface{}

	//获取data结构体信息
	keys, points, a1, a2 := getKeysAndPointers(data)

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

	return turnSliceBack(c, a1, a2)
}

func Get(c *gin.Context, table string, data interface{}, where map[string]interface{}) error {
	var wh []string
	var keys []string
	var values []interface{}

	//取字段名
	{
		g := reflect.ValueOf(data).Elem().Index(0)
		keys = getKeys(reflect.TypeOf(g.Interface()))
	}

	for k, v := range where {
		wh = append(wh, k+"=?")
		values = append(values, v)
	}
	var w string
	if len(where) != 0 {
		w = fmt.Sprintf("WHERE %s", strings.Join(wh, ","))
	}
	SQL := fmt.Sprintf("SELECT %s FROM %s %s", strings.Join(keys, ","), table, w)
	if rows, err := DB.Query(SQL, values...); err != nil {
		er(c, err)
		return err
	} else {
		for t := 1; rows.Next(); t++ {
			e := reflect.ValueOf(data).Elem().Index(0).Interface()
			_, g, a1, a2 := getKeysAndPointers(&e)
			if err := rows.Scan(g...); err != nil {
				er(c, err)
				return err
			}
			//数组还原
			if err := turnSliceBack(c, a1, a2); err != nil {
				return err
			}

			if reflect.ValueOf(data).Elem().Len() <= t {
				reflect.ValueOf(data).Elem().Index(t - 1).Set(reflect.ValueOf(e))
			} else {
				reflect.Append(reflect.ValueOf(data).Elem(), reflect.ValueOf(e))
			}
		}
	}
	return nil
}

func GetWithLimit(c *gin.Context, table string, data interface{}, where map[string]interface{}, page int) error { //传入data为结构体的切片的指针，必须make，按容量扫描
	var keys []string
	var wh []string
	var values []interface{}

	var limit int

	{ //取字段名
		limit = reflect.ValueOf(data).Elem().Len()
		j := reflect.TypeOf(reflect.ValueOf(data).Elem().Index(0).Interface())
		keys = getKeys(j)
	}

	for k, v := range where {
		wh = append(wh, k+"=?")
		values = append(values, v)
	}
	page = (page - 1) * limit
	var w string
	if len(wh) != 0 { //应对不需要where的情况
		w = "WHERE " + strings.Join(wh, ",")
	}
	SQL := fmt.Sprintf("SELECT %s FROM %s %s limit %d,%d", strings.Join(keys, ","), table, w, page, page+limit)
	rows, err := DB.Query(SQL, values...)
	defer rows.Close()
	if err != nil {
		er(c, err)
		return err
	}

	//数组json临时存储
	var a1 []*string
	var a2 []interface{}

	i := 0
	for rows.Next() {
		//取指针
		var points []interface{}
		o := reflect.ValueOf(data).Elem().Index(i).Interface()
		_, points, b1, b2 := getKeysAndPointers(&o)
		a1 = append(a1, b1...)
		a2 = append(a2, b2...)

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

	return turnSliceBack(c, a1, a2)
}

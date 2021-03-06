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
	if e != nil && c != nil {
		fmt.Println(e) //DEMO
		Modules.CallBack.Error(c, 102)
	}
}

func erCode(c *gin.Context, code int) {
	if c != nil {
		Modules.CallBack.Error(c, code)
	}
}

func databaseConfig()string{
	formatConfig:= func(data map[string]string) string{
		return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v",data["username"],data["password"],data["host"],data["port"],data["database"])
	}
	if !Modules.File.Exists("config.ini"){
		//引导链接数据库
		fmt.Println("[乌拉，没找到配置文件]Mysql链接向导开始")
		var data = make(map[string]string)
		var t string
		var need =[]string{
			"Host",
			"Port",
			"UserName",
			"Password",
			"Database",
		}
		for i:=0;i<len(need);i++{
			fmt.Printf(need[i]+":")
			_, _ = fmt.Scanln(&t)
			if t == ""{
				i--
				continue
			}
			data[strings.ToLower(need[i])]=t
		}
		file,_:=json.Marshal(data)
		save:err := Modules.File.Write("config.ini", file)
		if err==nil{
			fmt.Println("已保存为【config.ini】")
		}else{
			fmt.Println(err)
			fmt.Println("保存失败，是否重试")
			for{
				fmt.Printf("Y/N:")
				_, _ = fmt.Scanln(&t)
				t=strings.ToLower(t)
				if t=="y"{
					goto save
				}else if t=="n"{
					break
				}
			}
		}
		return formatConfig(data)
	}
	givUp:=func()string{//处理err
		fmt.Println("读取配置文件【config.ini】出错，请重新配置")
		if Modules.File.Remove("config.ini")!=nil{
			fmt.Println("无法对【config.ini】操作，请检查文件权限或手动删除")
			os.Exit(3)
		}
		return databaseConfig()
	}
	//读取配置
	file,err:=Modules.File.Read("config.ini")
	if err!=nil{
		return givUp()
	}
	var data map[string]string
	err=json.Unmarshal(file,&data)
	if err!=nil{
		return givUp()
	}
	return formatConfig(data)
}

func InitDatabase() { //连接数据库
	var err error
	if DB, err = sql.Open("mysql", databaseConfig()); err != nil || DB.Ping() != nil {
		fmt.Println(DB.Ping())
		fmt.Println(err)
		os.Exit(3)
	}
	fmt.Println("数据库已连接")
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
		//数组字符串化
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice:
			var temp string
			if reflect.ValueOf(v).Len() == 0 {
				temp = ""
			} else {
				t, _ := json.Marshal(v)
				temp = string(t)
			}
			//DEMO 应对超长
			values = append(values, temp)
		default:
			values = append(values, v)
		}
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
		//数组字符串化
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice:
			temp, _ := json.Marshal(v)
			//DEMO 应对超长
			values = append(values, temp)
		default:
			values = append(values, v)
		}
	}
	{
		var temp []interface{}
		wh, temp = Modules.Tool.MakeWhere(where)
		values = append(values, temp)
	}

	SQL := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(keys, ","), strings.Join(wh, ","))
	id, err := Exec(SQL, values...)
	er(c, err)
	return id, err
}

func Delete(c *gin.Context, table string, where map[string]interface{}) (int64, error) {
	var wh []string
	var values []interface{}
	wh, values = Modules.Tool.MakeWhere(where)
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
		keys = append(keys, strings.ToLower(v.Field(i).Name))
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
			a2 = append(a2, f.Field(i).Addr().Interface())
			continue
		}
		pointers = append(pointers, f.Field(i).Addr().Interface())
	}
	return keys, pointers, a1, a2
}

func turnSliceBack(c *gin.Context, a1 []*string, a2 []interface{}) error {
	for ii, v := range a1 {
		temp := reflect.New(reflect.TypeOf(a2[ii]).Elem())
		if *v != "" && json.Unmarshal([]byte(*v), temp.Interface()) != nil {
			err := errors.New("未知错误-解码失败")
			er(c, err)
			return err
		}
		if temp.Elem().Len() != 0 {
			reflect.ValueOf(a2[ii]).Elem().Set(temp.Elem())
		} else {
			//论如何空手初始化未知类型切片……
			reflect.ValueOf(a2[ii]).Elem().Set(reflect.MakeSlice(temp.Elem().Type(), 0, 0))
		}
	}
	return nil
}

func GetRow(c *gin.Context, table string, data interface{}, where map[string]interface{}) error { //传入的data为结构体的指针，否则无法工作
	var wh []string
	var values []interface{}

	//获取data结构体信息
	keys, points, a1, a2 := getKeysAndPointers(data)

	//构造查询语句
	wh, values = Modules.Tool.MakeWhere(where) //没有考虑不需要where，因为没用到

	SQL := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(keys, ","), table, strings.Join(wh, ","))
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

	//初始化data
	s := reflect.ValueOf(data).Elem()
	s.Set(reflect.MakeSlice(s.Type(), 1, 1))

	//取字段名
	{
		g := s.Index(0)
		keys = getKeys(reflect.TypeOf(g.Interface()))
	}

	wh, values = Modules.Tool.MakeWhere(where)
	var w string
	if where != nil && len(where) != 0 {
		w = fmt.Sprintf("WHERE %s", strings.Join(wh, ","))
	}
	SQL := fmt.Sprintf("SELECT %s FROM %s %s", strings.Join(keys, ","), table, w)
	if rows, err := DB.Query(SQL, values...); err != nil {
		rows.Close()
		er(c, err)
		return err
	} else {
		defer rows.Close()
		for t := 1; rows.Next(); t++ {
			e := reflect.New(s.Index(0).Type())
			_, g, a1, a2 := getKeysAndPointers(e.Interface())
			if err := rows.Scan(g...); err != nil {
				er(c, err)
				return err
			}

			//数组还原
			if err := turnSliceBack(c, a1, a2); err != nil {
				return err
			}

			if t == 1 {
				s.Index(t - 1).Set(e.Elem())
			} else {
				s.Set(reflect.Append(reflect.ValueOf(data).Elem(), e.Elem()))
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

	wh, values = Modules.Tool.MakeWhere(where)
	page = (page - 1) * limit
	var w string
	if where != nil && len(wh) != 0 { //应对不需要where的情况
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
		_, points, b1, b2 := getKeysAndPointers(reflect.ValueOf(data).Elem().Index(i).Addr().Interface())
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
		erCode(c, 115)
		return errors.New("没有更多了")
	} else if i != limit {
		reflect.ValueOf(data).Elem().SetLen(i)
	}

	return turnSliceBack(c, a1, a2)
}

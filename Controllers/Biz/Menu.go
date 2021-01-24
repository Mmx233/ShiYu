package Controllers

import (
	"Mmx/Modules"
	"Mmx/Service"
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
)

type menu struct {
	Fav fav
}

func get(c *gin.Context, where map[string]interface{}) (interface{}, error) {
	type d struct {
		Id         uint   `json:"id"`
		BizId      uint   `json:"biz_id"`
		Name       string `json:"name"`
		Price      string `json:"price"`
		CatFoodId  uint   `json:"cat_food_id"`
		IsFavorite bool   `json:"is_favorite" skip:"true"`
	}
	data := make([]d, 1)
	if err := Service.Get(c, "biz_menu", &data, where); err != nil {
		return nil, err
	}
	//favorite相关
	var u struct {
		Favorites []uint `json:"favorites"`
	}
	if err := Service.GetRow(c, "user", &u, map[string]interface{}{
		"username": c.GetString("username"),
	}); err != nil {
		return nil, err
	}
	for i := 0; i < len(data); i++ {
		//Go没有find，好烦啊
		data[i].IsFavorite = Modules.Tool.Find(u.Favorites, data[i].Id)
	}
	return data, nil
}

func (*menu) Information(c *gin.Context) {
	var form struct {
		Id uint `form:"id" binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if Service.Checker.BizExist(form.Id) {
		Modules.CallBack.Error(c, 116)
		return
	}
	data, err := get(c, map[string]interface{}{
		"biz_id": form.Id,
	})
	if err != nil {
		return
	}
	Modules.CallBack.Success(c, data)
}

func (*menu) InformationForCat(c *gin.Context) {
	var form struct {
		Id uint `form:"id" binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Service.Checker.CatIdExist("food", form.Id) {
		Modules.CallBack.Error(c, 119)
		return
	}
	data, err := get(c, map[string]interface{}{
		"cat_food_id": form.Id,
	})
	if err != nil {
		return
	}
	Modules.CallBack.Success(c, data)
}

func (*menu) InformationForFavorites(c *gin.Context) {
	data, err := get(c, nil)
	if err != nil {
		return
	}
	var u struct {
		Favorites []uint `json:"favorites"`
	}
	if Service.GetRow(c, "user", &u, map[string]interface{}{
		"username": c.GetString("username"),
	}) != nil {
		return
	}
	d := reflect.ValueOf(data)
	var favorites []interface{}
	for i := 0; i < d.Len(); i++ {
		if Modules.Tool.Find(u.Favorites, uint(d.Index(i).FieldByName("Id").Uint())) {
			favorites = append(favorites, d.Field(i).Interface())
		}
	}
	Modules.CallBack.Success(c, favorites)
}

func (*menu) New(c *gin.Context) {
	var form struct {
		BizId     uint    `form:"biz_id" binding:"required"`
		CatFoodId uint    `form:"cat_food_id" binding:"required"`
		Name      string  `form:"name" binding:"required"`
		Price     float32 `form:"price" binding:"min=0,max=999"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Modules.Checker.Form(c, &form) {
		return
	}
	if !Service.Checker.BizExist(form.BizId) {
		Modules.CallBack.Error(c, 116)
		return
	}
	if !Service.Checker.CatIdExist("food", form.CatFoodId) {
		Modules.CallBack.Error(c, 119)
		return
	}

	id, err := Service.Insert(c, "biz_menu", map[string]interface{}{
		"biz_id":      form.BizId,
		"name":        form.Name,
		"price":       form.Price,
		"cat_food_id": form.CatFoodId,
	})
	if err != nil {
		return
	}
	Modules.CallBack.Success(c, map[string]interface{}{
		"id": id,
	})
}

func (*menu) Renew(c *gin.Context) {
	var form struct {
		Id        uint    `form:"id" binding:"required"`
		BizId     uint    `form:"biz_id" binding:"required"`
		CatFoodId uint    `form:"cat_food_id" binding:"required"`
		Name      string  `form:"name" binding:"required"`
		Price     float32 `form:"price" binding:"min=0,max=999"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Modules.Checker.Form(c, &form) {
		return
	}
	if !Service.Checker.MenuExist(form.Id) {
		Modules.CallBack.Error(c, 123)
		return
	}
	if !Service.Checker.BizExist(form.BizId) {
		Modules.CallBack.Error(c, 116)
		return
	}
	if !Service.Checker.CatIdExist("food", form.CatFoodId) {
		Modules.CallBack.Error(c, 119)
		return
	}
	if _, err := Service.Update(c, "biz_menu", map[string]interface{}{
		"biz_id":      form.BizId,
		"name":        form.Name,
		"price":       form.Price,
		"cat_food_id": form.CatFoodId,
	}, map[string]interface{}{
		"id": form.Id,
	}); err != nil {
		return
	}
	Modules.CallBack.Default(c)
}

func (*menu) Change(c *gin.Context) {
	var form struct {
		Id     uint        `form:"id" binding:"required"`
		Target string      `form:"target" binding:"required"`
		Value  interface{} `form:"value"  binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Service.Checker.MenuExist(form.Id) {
		Modules.CallBack.Error(c, 123)
		return
	}
	form.Target = strings.ToLower(form.Target)
	switch form.Target {
	case "biz_id":
		if _, ok := form.Value.(uint); !ok {
			Modules.CallBack.Error(c, 101)
			return
		}
		if !Service.Checker.BizExist(form.Value.(uint)) {
			Modules.CallBack.Error(c, 116)
			return
		}
	case "name":
		if _, ok := form.Value.(string); !ok {
			Modules.CallBack.Error(c, 101)
			return
		}
		if !Modules.Checker.Name(c, form.Value.(string)) {
			return
		}
	case "price":
		if _, ok := form.Value.(float32); !ok {
			Modules.CallBack.Error(c, 101)
			return
		}
		if form.Value.(float32) < 0 {
			Modules.CallBack.ErrorWithErr(c, 101, errors.New("价格不能为负"))
			return
		}
	case "cat_food_id":
		if _, ok := form.Value.(uint); !ok {
			Modules.CallBack.Error(c, 101)
			return
		}
		if !Service.Checker.CatIdExist("food", form.Value.(uint)) {
			return
		}
	default:
		Modules.CallBack.Error(c, 114)
		return
	}
	if _, err := Service.Update(c, "biz_menu", map[string]interface{}{
		form.Target: form.Value,
	}, map[string]interface{}{
		"id": form.Id,
	}); err != nil {
		return
	}
	Modules.CallBack.Default(c)
}

func (*menu) Delete(c *gin.Context) {
	var form struct {
		Id uint `form:"id" binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Service.Checker.MenuExist(form.Id) {
		Modules.CallBack.Error(c, 123)
		return
	}
	if _, err := Service.Delete(c, "biz_menu", map[string]interface{}{
		"id": form.Id,
	}); err != nil {
		return
	}
	Modules.CallBack.Default(c)
}

type fav struct{}

func (*fav) Make(c *gin.Context) {
	var form struct {
		Id uint `form:"id" binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	if !Service.Checker.MenuExist(form.Id) {
		Modules.CallBack.Error(c, 123)
		return
	}
	var fav []uint
	var is bool
	if is, fav = Service.Checker.IsFav(c, form.Id); is {
		Modules.CallBack.Error(c, 125)
		return
	}
	fav = append(fav, form.Id)
	if _, err := Service.Update(c, "user", map[string]interface{}{
		"favorites": fav,
	}, map[string]interface{}{
		"username": c.GetString("username"),
	}); err != nil {
		return
	}
	Modules.CallBack.Default(c)
}

func (*fav) Cancel(c *gin.Context) {
	var form struct {
		Id uint `form:"id" binding:"required"`
	}
	if !Modules.Tool.BindForm(c, &form) {
		return
	}
	var fav []uint
	var is bool
	if is, fav = Service.Checker.IsFav(c, form.Id); !is {
		Modules.CallBack.Error(c, 124)
		return
	}
	for i, v := range fav {
		if v == form.Id {
			if i == len(fav)-1 {
				if i != 0 {
					fav = fav[:len(fav)-2]
				} else {
					fav = make([]uint, 0)
				}
			} else if i == 0 {
				fav = fav[1:]
			} else {
				fav = append(fav[:i-1], fav[i+1:]...)
			}
			break
		}
	}
	if _, err := Service.Update(c, "user", map[string]interface{}{
		"favorites": fav,
	}, map[string]interface{}{
		"username": c.GetString("username"),
	}); err != nil {
		return
	}
	Modules.CallBack.Default(c)
}

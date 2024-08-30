package test

import (
	"github.com/mosongcc/gotool/gjson"
	"github.com/mosongcc/gotool/gsql"
	_ "modernc.org/sqlite"
	"reflect"
	"testing"
)

// _ "modernc.org/sqlite"  go写的sqlite实现，兼容性更好。

type User struct {
	Name  string
	Age   uint
	Ctime string
}

func Test_CURD(t *testing.T) {
	var err error
	defer func() {
		if err != nil {
			t.Error(err.Error())
		}
	}()

	var u = &User{}

	//打开连接
	db, err := gsql.Open(gsql.Sqlite, "gsql.sqlite.db")
	if err != nil {
		return
	}

	//新增数据
	_, err = db.Insert(&User{Name: "喵喵123", Age: 10, Ctime: "20240101"}).Exec()
	if err != nil {
		return
	}

	//查询数据
	userList, err := gsql.Find[User](db.Select(u, &u.Name, &u.Ctime).Where(&u.Name, gsql.Eq, "喵喵123").Limit(0, 10).Query())
	if err != nil {
		return
	}
	t.Log(len(userList), "   ", gjson.MarshalString(userList))

	//修改数据
	_, err = db.Update(u, map[any]any{&u.Name: "喵喵321", &u.Ctime: "20240201"}).Where(&u.Name, gsql.Eq, "喵喵123").Exec()
	if err != nil {
		return
	}

	//查询数据
	userList, err = gsql.Find[User](db.Select(u, &u.Name, &u.Ctime).Where(&u.Name, gsql.Eq, "喵喵321").Limit(0, 10).Query())
	if err != nil {
		return
	}
	t.Log(len(userList), "   ", gjson.MarshalString(userList))

	//删除数据
	_, err = db.Delete(u).Where(&u.Name, gsql.Eq, "喵喵").Exec()
	if err != nil {
		return
	}

}

func TestS(t *testing.T) {
	var u *User
	t.Log(reflect.TypeOf(u).String())
	t.Log(reflect.ValueOf(u).Interface())
	t.Log(reflect.ValueOf(&u).Kind())

	var u2 User
	t.Log(reflect.TypeOf(u2).String())
	t.Log(reflect.ValueOf(u2).Interface())
	t.Log(reflect.ValueOf(u2).Kind())

}

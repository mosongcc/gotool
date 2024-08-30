package test

import (
	"github.com/mosongcc/gotool/gjson"
	"github.com/mosongcc/gotool/gsql"
	_ "modernc.org/sqlite"
	"testing"
)

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
	_, err = db.Insert(&User{Name: "喵喵", Age: 10, Ctime: "20240101"}).Exec()
	if err != nil {
		return
	}

	//查询数据
	userList, err := gsql.Find[User](db.Select(u, &u.Name, &u.Ctime).Where(&u.Ctime, gsql.Gt, "20230101").Limit(0, 10).Query())
	if err != nil {
		return
	}
	t.Log(len(userList), "   ", gjson.MarshalString(userList))

	//修改数据
	_, err = db.Update(u, map[any]any{&u.Name: "旺旺", &u.Ctime: "20240201"}).Where(&u.Name, gsql.Eq, "喵喵").And(&u.Ctime, gsql.Gt, "20240101").Exec()
	if err != nil {
		return
	}

	//删除数据
	_, err = db.Delete(u).Where(&u.Name, gsql.Eq, "旺旺").Exec()
	if err != nil {
		return
	}

}

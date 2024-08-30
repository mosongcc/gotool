package gsql

import (
	"github.com/mosongcc/gotool/gjson"
	"reflect"
	"testing"
)

func Test_getTableName(t *testing.T) {
	type CatColor struct {
		Name string
		Age  uint64
	}
	var cat = CatColor{}

	typeOf := reflect.TypeOf(cat)
	valueOf := reflect.ValueOf(cat)

	name := getTableName(typeOf, valueOf)
	t.Log(name)
}

func TestTN(t *testing.T) {
	type CatColor struct {
		CatName string
		Age     uint64
	}
	var cat = &CatColor{}
	name := TN(cat)
	t.Log("TableName:", name)
	t.Log(FN(&cat.CatName), "  ", FN(&cat.Age))
}

func Test_getStructFields(t *testing.T) {
	type CatColor struct {
		NameCn string
		Age    uint64
	}
	var cat = &CatColor{NameCn: "猫名"}

	typeOf := reflect.TypeOf(cat)
	valueOf := reflect.ValueOf(cat)

	k, p, a := getStructFields(typeOf, valueOf)
	t.Log(gjson.MarshalString(k))
	t.Log(gjson.MarshalString(p))
	t.Log(gjson.MarshalString(a))
}

package gsql

import (
	"log/slog"
	"reflect"
	"testing"
	"time"
)

type Cat struct {
	Name string
	Age  int64
}

func (Cat) TableName() string {
	return "cat"
}

func Test_getTableName(t *testing.T) {
	t.Log(getTableName(&Cat{}))

	beg := time.Now()
	num := 100000
	for i := 0; i < num; i++ {
		getTableName(&Cat{})
	}
	t.Logf("执行%v次，耗时：%v毫秒", num, time.Since(beg).Milliseconds())
	//执行100000次，耗时：39毫秒
}

func Test_re(t *testing.T) {

	cat := &Cat{}

	catType := reflect.TypeOf(cat)
	//catValue := reflect.ValueOf(cat)

	slog.Info(catType.String()) //*gsql.Cat

	if method, ok := catType.MethodByName("TableName"); ok {
		slog.Info(method.Name)
		name := method.Func.Call([]reflect.Value{reflect.ValueOf(cat)})[0].String()
		slog.Info(name)
	}

	t.Logf("%v", getTableName(cat))
	t.Logf("%v", getFieldName(&cat.Name))
}

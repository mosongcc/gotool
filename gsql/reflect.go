package gsql

import (
	"fmt"
	"github.com/mosongcc/gotool/gstring"
	"reflect"
	"strings"
	"sync"
)

// 缓存表名与字段名 key=指针地址 value=名字字符串
var ptrMap sync.Map

func setPtrMap(ptr uintptr, v string) {
	ptrMap.Store(ptr, v)
}

func getPtrMap(ptr uintptr) string {
	if v, ok := ptrMap.Load(ptr); ok {
		return v.(string)
	}
	return ""
}

// TN 根据传入的表信息，获取表名
func TN(table any) string {
	valueOf := reflect.ValueOf(table)
	if reflect.Pointer == valueOf.Kind() {
		return getTableNameByCache(reflect.TypeOf(table), valueOf)
	} else {
		return fmt.Sprintf("%v", table)
	}
}

// getTableNameByCache 反射表名,优先从TableName方法获取,没有方法则从名字获取
func getTableNameByCache(typeOf reflect.Type, valueOf reflect.Value) (name string) {

	name = getPtrMap(valueOf.Pointer())
	if name != "" {
		return
	}
	name = getTableName(typeOf, valueOf)

	//缓存表名字
	setPtrMap(valueOf.Pointer(), name)

	//缓存字段名
	for j := 0; j < valueOf.Elem().NumField(); j++ {
		fieldPointer := valueOf.Elem().Field(j).Addr().Pointer()
		fieldName := gstring.Underline(typeOf.Elem().Field(j).Name)
		setPtrMap(fieldPointer, fieldName)
	}

	return
}

func getTableName(typeOf reflect.Type, valueOf reflect.Value) (name string) {
	// 优先函数取表名
	method, isSet := typeOf.MethodByName("TableName")
	if isSet {
		res := method.Func.Call([]reflect.Value{valueOf})
		name = res[0].String()
	} else {
		slices := strings.Split(typeOf.String(), ".")
		name = gstring.Underline(slices[len(slices)-1])
	}
	return
}

// FN 获取结构体字段名
func FN(field any) string {
	valueOf := reflect.ValueOf(field)
	if reflect.Pointer == valueOf.Kind() {
		// 注意：从缓存取字段名，必须先获取表名
		return getPtrMap(valueOf.Pointer())
	} else {
		return fmt.Sprintf("%v", valueOf)
	}
}

// 取结构体字段
func getStructFields(typeOf reflect.Type, valueOf reflect.Value) (keys, place []string, args []any) {
	for i := 0; i < valueOf.Elem().NumField(); i++ {
		keys = append(keys, gstring.Underline(typeOf.Elem().Field(i).Name))
		place = append(place, "?")
		args = append(args, valueOf.Elem().Field(i).Field(0).Field(0).Interface())
	}
	return
}

package gsql

import (
	"fmt"
	"github.com/mosongcc/gotool/gstring"
	"reflect"
	"strings"
)

// GetTN 根据传入的表信息，获取表名
func GetTN(table any) string {
	valueOf := reflect.ValueOf(table)
	if reflect.Pointer == valueOf.Kind() {
		return GetTNByReflect(reflect.TypeOf(table), valueOf)
	} else {
		return fmt.Sprintf("%v", table)
	}
}

// GetTNByReflect 反射表名,优先从TableName方法获取,没有方法则从名字获取
func GetTNByReflect(typeOf reflect.Type, valueOf reflect.Value) (name string) {
	method, isSet := typeOf.MethodByName("TableName")
	if isSet {
		res := method.Func.Call([]reflect.Value{valueOf})
		name = res[0].String()
	} else {
		slices := strings.Split(typeOf.String(), ".")
		name = gstring.Underline(slices[len(slices)-1])
	}
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

// GetFN 获取结构体字段名
func GetFN(field any) string {
	valueOf := reflect.ValueOf(field)
	if reflect.Pointer == valueOf.Kind() {
		// 注意：取字段名之前，请先获取表名
		return getPtrMap(valueOf.Pointer())
	} else {
		return fmt.Sprintf("%v", valueOf)
	}
}

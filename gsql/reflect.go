package gsql

import (
	"fmt"
	"github.com/mosongcc/gotool/gstring"
	"reflect"
	"strings"
	"sync"
)

var (
	tableMap sync.Map //缓存表名
	fieldMap sync.Map //缓存字段名

)

func setTableMap(ptr uintptr, v string) {
	tableMap.Store(ptr, v)
}

func getTableMap(ptr uintptr) string {
	if v, ok := tableMap.Load(ptr); ok {
		return v.(string)
	}
	return ""
}

func setFieldMap(ptr uintptr, v string) {
	fieldMap.Store(ptr, v)
}

func getFieldMap(ptr uintptr) string {
	if v, ok := fieldMap.Load(ptr); ok {
		return v.(string)
	}
	return ""
}

// TN 根据传入的表信息，获取表名。  参数必须是结构体指针地址
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
	vptr := valueOf.Pointer()
	name = getTableMap(vptr)
	if name != "" {
		return
	}
	name = getTableName(typeOf, valueOf)
	setTableMap(vptr, name) //缓存表名字

	for j := 0; j < valueOf.Elem().NumField(); j++ {
		fieldPointer := valueOf.Elem().Field(j).Addr().Pointer()
		fieldName := gstring.Underline(typeOf.Elem().Field(j).Name)
		setFieldMap(fieldPointer, fieldName) //缓存字段名
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
		return getFieldMap(valueOf.Pointer())
	} else {
		return fmt.Sprintf("%v", valueOf)
	}
}

// 取结构体字段
func getStructFields(typeOf reflect.Type, valueOf reflect.Value) (keys, place []string, args []any) {
	if reflect.Pointer == valueOf.Kind() {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	for i := 0; i < typeOf.NumField(); i++ {
		name := typeOf.Field(i).Name
		value := valueOf.Field(i).Interface()

		keys = append(keys, gstring.Underline(name))
		place = append(place, "?")
		args = append(args, value)
	}
	return
}

package gsql

import "sync"

// 缓存表名与字段名 key=指针地址 value=名字字符串
// 取表名的同时缓存字段名
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

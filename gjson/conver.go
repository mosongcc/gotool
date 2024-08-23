package gjson

import "encoding/json"

// Convert 对象通过json中转转换
func Convert[S any, T any](source S) (target T, err error) {
	b := Marshal(source)
	err = json.Unmarshal(b, &target)
	if err != nil {
		return
	}
	return
}

// MustConvert 对象通过json中转转换
func MustConvert[S any, T any](source S) (target T) {
	return Must(Convert[S, T](source))
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

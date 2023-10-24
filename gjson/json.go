package gjson

import (
	"bytes"
	"encoding/json"
)

// Stringify 对象转为json字符串
func Stringify(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Parse 字符串解析为JSON对象
func Parse[T any](s string) (r T, err error) {
	err = json.Unmarshal([]byte(s), &r)
	if err != nil {
		return
	}
	return
}

// Convert 对象通过json中转转换
func Convert[S any, T any](source S) (target T, err error) {
	b, _ := json.Marshal(source)
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

// Must 忽略异常
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// Marshal  escapeHTML=false 则不转义符号 & < >
func Marshal(v any, escapeHTML bool) (b []byte, err error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(escapeHTML)
	err = jsonEncoder.Encode(v)
	b = bf.Bytes()
	return
}

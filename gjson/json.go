package gjson

import (
	"bytes"
	"encoding/json"
)

// Marshal 对象转为JSON
func Marshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// MarshalString 转为JSON字符串
func MarshalString(v any) string {
	return string(Marshal(v))
}

// Unmarshal 解析为JSON对象
func Unmarshal[T any](v []byte) (r T, err error) {
	err = json.Unmarshal(v, &r)
	if err != nil {
		return
	}
	return
}

// UnmarshalString 字符串解析为JSON对象
func UnmarshalString[T any](v string) (r T, err error) {
	return Unmarshal[T]([]byte(v))
}

// MarshalEscapeHTML  escapeHTML=false 则不转义符号 & < >
func MarshalEscapeHTML(v any, escapeHTML bool) (b []byte, err error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(escapeHTML)
	err = jsonEncoder.Encode(v)
	b = bf.Bytes()
	return
}

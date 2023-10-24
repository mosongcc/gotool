package gtype

import (
	"errors"
	"strconv"
	"strings"
)

// Int64 自定义int64类型，解析json时，支持数字与字符串格式 ，用于接收json参数解析时，第三方给的数字类型参数带了双引号
// 注意：当字段为空时，默认值是 0
type Int64 int64

func (t Int64) String() string {
	return strconv.FormatInt(int64(t), 10)
}

func (t Int64) MarshalJSON() ([]byte, error) {
	b := []byte(t.String())
	return b, nil
}

func (t *Int64) UnmarshalJSON(data []byte) (err error) {
	v := strings.Trim(string(data), "\"")
	if data == nil || v == "" || v == "null" {
		return
	}
	pint, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		err = errors.New("参数【" + string(data) + "】非预期的int64类型")
		return
	}
	*t = Int64(pint)
	return
}

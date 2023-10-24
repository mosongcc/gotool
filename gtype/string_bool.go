package gtype

import (
	"strings"
)

// Bool 字符串与数字类型的0与1都可直接使用
type Bool string

const (
	True  Bool = "1"
	False Bool = "0"
)

func (t *Bool) MarshalJSON() ([]byte, error) {
	return []byte("\"" + *t + "\""), nil
}

// UnmarshalJSON json解析时支持字符串与数字 1 与 0
func (t *Bool) UnmarshalJSON(data []byte) (err error) {
	if data == nil {
		return
	}
	*t = Bool(strings.Trim(string(data), "\""))
	return
}

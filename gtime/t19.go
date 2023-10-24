package gtime

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	LayoutT19 = "2006-01-02 15:04:05"
)

type T19 time.Time

func (nt T19) String() string {
	return time.Time(nt).Local().Format(LayoutT19)
}

func (nt T19) MarshalJSON() ([]byte, error) {
	t := time.Time(nt).Local()

	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(LayoutT19)+2)
	b = append(b, '"')
	if !t.IsZero() {
		b = t.AppendFormat(b, LayoutT19)
	}
	b = append(b, '"')
	return b, nil
}

func (nt *T19) UnmarshalJSON(data []byte) error {
	v := strings.Trim(string(data), "\"")
	if data == nil || v == "" {
		*nt = T19(time.Time{})
		return nil
	}
	t, err := time.Parse(LayoutT19, v)
	if err != nil {
		return fmt.Errorf("解析时间字符串'%s'出错", v)
	}
	*nt = T19(t.Local())
	return err
}

// NowF19 取当前格式化时间 如："2006-01-02 15:04:05"
func NowF19() string {
	return T19(time.Now()).String()
}

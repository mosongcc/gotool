package gtime

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	LayoutT8 = "20060102"
)

type T8 time.Time

func (nt T8) String() string {
	return time.Time(nt).Local().Format(LayoutT8)
}

func (nt T8) MarshalJSON() ([]byte, error) {
	t := time.Time(nt).Local()

	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(LayoutT8)+2)
	b = append(b, '"')
	if !t.IsZero() {
		b = t.AppendFormat(b, LayoutT8)
	}
	b = append(b, '"')
	return b, nil
}

func (nt *T8) UnmarshalJSON(data []byte) error {
	v := strings.Trim(string(data), "\"")
	if data == nil || v == "" {
		*nt = T8(time.Time{})
		return nil
	}
	t, err := time.Parse(LayoutT8, v)
	if err != nil {
		return fmt.Errorf("解析时间字符串'%s'出错", v)
	}
	*nt = T8(t.Local())
	return err
}

func NowF8() string {
	return T8(time.Now()).String()
}

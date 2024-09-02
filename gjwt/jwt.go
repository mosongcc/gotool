package gjwt

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var ErrToken = errors.New("token authentication failed")

// Encode 编码
func Encode(v any, exp time.Time, secret string) (token string) {
	secret = fmt.Sprintf("%x", md5.Sum([]byte(secret)))
	if exp.IsZero() {
		exp = time.Now().Add(15 * time.Hour)
	}
	valBytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	var vMap = map[string]any{"v": string(valBytes), "t": exp}

	b, err := json.Marshal(vMap)
	if err != nil {
		panic(err)
	}

	encodeValue, err := aesEncodeCBC(string(b), secret)
	if err != nil {
		panic(err)
		return
	}
	token = sign(encodeValue + secret)[:12] + encodeValue
	return
}

// Decode 解码
func Decode[T any](token string, secret string) (v T, err error) {
	secret = fmt.Sprintf("%x", md5.Sum([]byte(secret)))
	if len(token) < 32 {
		err = ErrToken
		return
	}
	tokenVal := token[12:]
	if token[:12] != sign(tokenVal + secret)[:12] {
		err = ErrToken
		return
	}
	tokenVal, err = aesDecodeCBC(tokenVal, secret)
	if err != nil {
		return
	}

	var vMap map[string]any
	err = json.Unmarshal([]byte(tokenVal), &vMap)
	if err != nil {
		return
	}

	// 判断有效期
	exp := vMap["t"].(time.Time)
	if time.Now().After(exp) {
		err = ErrToken
		return
	}

	// 取值
	val := vMap["v"].(string)
	err = json.Unmarshal([]byte(val), &v)
	if err != nil {
		return
	}
	return
}

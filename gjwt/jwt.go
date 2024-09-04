package gjwt

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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

	var vMap = map[string]string{"v": string(valBytes), "t": exp.Format(time.RFC3339)}

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
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
		if err != nil {
			slog.Info("JWT Decode : " + err.Error())
			err = ErrToken
			return
		}
	}()

	secret = fmt.Sprintf("%x", md5.Sum([]byte(secret)))
	if len(token) < 32 {
		err = errors.New("invalid token, token too short")
		return
	}
	tokenVal := token[12:]
	if token[:12] != sign(tokenVal + secret)[:12] {
		err = errors.New("invalid token")
		return
	}
	tokenVal, err = aesDecodeCBC(tokenVal, secret)
	if err != nil {
		return
	}

	var vMap map[string]string
	err = json.Unmarshal([]byte(tokenVal), &vMap)
	if err != nil || vMap == nil {
		return
	}

	// 判断有效期
	expV, ok := vMap["t"]
	if !ok || len(expV) != 14 {
		err = errors.New("invalid token exp")
		return
	}
	exp, err := time.Parse(time.RFC3339, expV)
	if err != nil {
		return
	}
	if time.Now().After(exp) {
		err = errors.New("token expired")
		return
	}

	// 取值
	val := vMap["v"]
	err = json.Unmarshal([]byte(val), &v)
	if err != nil {
		return
	}
	return
}

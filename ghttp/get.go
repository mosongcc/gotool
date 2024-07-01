package ghttp

import (
	"io"
	"net/http"
)

// GetBody GET请求返回Body内容
func GetBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

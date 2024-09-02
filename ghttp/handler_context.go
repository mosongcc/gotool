package ghttp

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type ServeMux struct {
	*http.ServeMux
}

func NewServeMux() *ServeMux {
	return &ServeMux{ServeMux: http.NewServeMux()}
}

func (mux *ServeMux) GET(path string, handler HandleFunc) {
	mux.Handle("GET "+path, handler)
}

func (mux *ServeMux) POST(path string, handler HandleFunc) {
	mux.Handle("POST "+path, handler)
}

type Context struct {
	w http.ResponseWriter
	r *http.Request
	c context.Context
	e []error
}

// HandleFunc 自定义Handler处理http请求
type HandleFunc func(*Context) error

func (handler HandleFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ctx = &Context{c: context.Background(), w: w, r: r}
	defer func(ctx *Context) {
		if e := recover(); e != nil {
			ctx.Error(e.(error))
		}
		ctx.logger()
		ctx.resError()
	}(ctx)
	ctx.WithValue("log_req_path", ctx.r.URL.Path)
	ctx.Error(handler(ctx))
}

func (ctx *Context) resError() {
	var errMsg = ctx.errorMsg()
	if errMsg != "" {
		contentType := strings.ToLower(ctx.r.Header.Get("Content-Type"))

		//根据请求类型响应不同格式内容
		if strings.Contains(contentType, "application/json") {
			ctx.JSON(map[string]any{"error": errMsg})
		} else {
			ctx.TEXT(fmt.Sprintf("error: %s", errMsg))
		}
	}
}

func (ctx *Context) WithValue(k any, v any) {
	ctx.c = context.WithValue(ctx.c, k, v)
}

func (ctx *Context) Error(e error) {
	if e == nil {
		return
	}
	ctx.e = append(ctx.e, e)
}

func (ctx *Context) errorMsg() string {
	var errMsg = ""
	if ctx.e != nil && len(ctx.e) > 0 {
		for _, err := range ctx.e {
			if err == nil {
				continue
			}
			errMsg += err.Error() + "  |  "
		}
		// 防止日志太长
		if len(errMsg) > 10000 {
			errMsg = errMsg[:10000]
		}
	}
	return errMsg
}

// GetBody 取Body
func (ctx *Context) GetBody() (b []byte) {
	b, err := io.ReadAll(ctx.r.Body)
	if err != nil {
		panic(err)
	}
	ctx.WithValue("log_req_body", string(b))
	return
}

// Bind 绑定参数
func (ctx *Context) Bind(v any) {
	if err := json.Unmarshal(ctx.GetBody(), &v); err != nil {
		panic(err)
	}
	return
}

func Bind[T any](ctx *Context) (v T) {
	ctx.Bind(&v)
	return
}

// JSON 响应JSON
func (ctx *Context) JSON(data any) {
	ctx.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	ctx.WithValue("log_res_body", string(b))
	_, err = ctx.w.Write(b)
	if err != nil {
		panic(err)
	}
}

// HTML 渲染HTML
func (ctx *Context) HTML(tpl *template.Template, name string, data any) {
	ctx.w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(ctx.w, name, data); err != nil {
		panic(fmt.Errorf("RenderHTML Error name：%s  errMsg：%s", name, err.Error()))
	}
}

// TEXT 渲染TEXT
func (ctx *Context) TEXT(v string) {
	_, _ = ctx.w.Write([]byte(v))
	ctx.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

// NotFound 404
func (ctx *Context) NotFound() {
	_, _ = ctx.w.Write([]byte("404 Not Found"))
	ctx.w.WriteHeader(http.StatusNotFound)
}

// Redirect 302
func (ctx *Context) Redirect(url string) {
	http.Redirect(ctx.w, ctx.r, url, http.StatusFound)
}

// 输出接口请求响应日志
func (ctx *Context) logger() {
	var errMsg = ctx.errorMsg()
	if errMsg != "" {
		errMsg = "\n错误信息：" + errMsg
	}

	// json请求日志
	if strings.Contains(ctx.r.Header.Get("Content-Type"), "application/json") {
		slog.Info(fmt.Sprintf("APILOG \n接口地址：%s \n请求报文：%s \n响应报文 %s %s",
			ctx.c.Value("log_req_path"),
			ctx.c.Value("log_req_body"),
			ctx.c.Value("log_res_body"),
			errMsg))
		return
	}

}

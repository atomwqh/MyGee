package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// handlerFunc希望访问到解析的参数
// context中添加一个对象和方法来提供对路由参数的访问，并将解析后的参数存储到Params中
// Day_Five 添加中间件，中间件可以作用在处理流程前，也可以作用在处理流程后，为此我们继续修改Context

type H map[string]interface{} // 字符串——空接口：gee.H的别名

type Context struct {
	// Origin
	Writer http.ResponseWriter
	Req    *http.Request
	// requst
	Path   string
	Method string
	Params map[string]string
	// response
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int // 记录当前执行到第几个中间件
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// 创建了一个New结构体，接收两个参数分别代表HTTP响应写入器和HTTP请求，同时也将请求的路径和方法赋值给了Path和Method
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// Query 和 Postman参数的访问方法
func (c *Context) Postman(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// status + header 的展示
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 构造string，html，json格式的响应方法
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

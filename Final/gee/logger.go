package gee

import (
	"log"
	"time"
)

// 现在我们来设计中间件
// logger:记录请求到响应所花费的时间

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

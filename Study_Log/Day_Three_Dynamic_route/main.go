package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	// handler参数变成了gee.Context，添加了查询Query和Postman参数的函数
	// gee.Context 封装了string，json，html，快速构造http响应
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Day_two of myGee</h1>")
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		// dynamic
		// expect /hello/wqh
		c.String(http.StatusOK, "hello %s, you are reading %s\n", c.Param("name"), c.Path)
	})

	r.GET("/hello", func(c *gee.Context) {
		// 静态路由
		// expect /hello?name=wqh
		c.String(http.StatusOK, "hello %s, you are reading %s\n", c.Query("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *gee.Context) {
		// expect /assets/wqh/nb
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	r.Run(":6666")
}

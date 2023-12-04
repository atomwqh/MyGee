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
	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you are reading %s\n", c.Query("name"), c.Path)
	})
	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.Postman("username"),
			"password": c.Postman("password"),
		})
	})
	r.Run(":6666")
}

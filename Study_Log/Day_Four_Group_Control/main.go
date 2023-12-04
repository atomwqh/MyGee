package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	// 现在来测试分组路由
	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Index Pages</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})
		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you are at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.Postman("username"),
				"password": c.Postman("password"),
			})
		})
	}

	r.Run(":6666")
}

// 分组测试成功

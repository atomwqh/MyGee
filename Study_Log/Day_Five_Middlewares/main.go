package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func onlyForv2() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		c.Fail(500, "Inter Server Error")
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
func main() {
	r := gee.New()
	r.Use(gee.Logger()) // 全局中间件
	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Index Pages</h1>")
	})
	//v1 := r.Group("/v1")
	//{
	//	v1.GET("/", func(c *gee.Context) {
	//		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	//	})
	//	v1.GET("/hello", func(c *gee.Context) {
	//		c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"), c.Path)
	//	})
	//}
	v2 := r.Group("/v2")
	v2.Use(onlyForv2()) // 局部v2中间件
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you are at %s\n", c.Param("name"), c.Path)
		})
		//v2.POST("/login", func(c *gee.Context) {
		//	c.JSON(http.StatusOK, gee.H{
		//		"username": c.Postman("username"),
		//		"password": c.Postman("password"),
		//	})
		//})
	}

	r.Run(":6666")
}

// 分组测试成功

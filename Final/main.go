package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "good morning wqh\n")
	})
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"wqh"}
		c.String(http.StatusOK, names[100]) // 数组越界导致panic
	})
	r.Run(":9999")
}

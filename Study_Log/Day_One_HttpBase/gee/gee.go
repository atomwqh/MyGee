package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc 定义gee使用的请求接口
type HandlerFunc func(http.ResponseWriter, *http.Request)

// 提供一个路由映射表
type Engine struct {
	router map[string]HandlerFunc
}

// New函数创建一个gee.Enige
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 添加路由，将路由和后端函数绑定
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 解析请求的路径，查找路由映射表，如果查到，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprint(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

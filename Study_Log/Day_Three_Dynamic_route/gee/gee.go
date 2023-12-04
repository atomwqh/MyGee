package gee

import (
	"net/http"
)

// HandlerFunc 定义gee使用的请求接口
type HandlerFunc func(c *Context)

// 提供一个路由映射表
type Engine struct {
	router *router
}

// New函数创建一个gee.Enige
func New() *Engine {
	return &Engine{router: newRouter()}
}

// 添加路由，将路由和后端函数绑定
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
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
	c := newContext(w, req)
	engine.router.handle(c)
}

// 第二天后gee.go的代码精简了不少（抽离了route相关代码）
//最重要的还是实现了关于serveHTTP的接口，接管了所有http请求
//下一步的计划——继续构造Context对象

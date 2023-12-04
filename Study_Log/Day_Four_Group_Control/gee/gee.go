package gee

import (
	"log"
	"net/http"
)

// HandlerFunc 定义gee使用的请求接口
type HandlerFunc func(c *Context)

// 定义路由组的功能，由此之后和路由相关的函数都由RouteGroup来实现
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 支持中间件
	parent      *RouterGroup  // 支持嵌套
	engine      *Engine       // 所有group共享一个enige实例
}

// 提供一个路由映射表
// 这里将Enige作为最顶层的分组，具有RouterGroup的所有功能
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // 保存所有groups
}

// New函数创建一个gee.Enige
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 用于创建新的路由组
// 这里所有的路由组都共享一个Enige实例
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix,
		middlewares: nil,
		parent:      group,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 之后的相关路由操作都由group来接手
// 添加路由，将路由和后端函数绑定
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp // 这里之后就缺少前缀了，之后就要加上分组后约定的前缀
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
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

// day_four 开始实现分组控制
/*
这里简单介绍分组
针对每一个路由进行控制
like /post 匿名可访问
/admin 需要鉴权访问
/api 需要RSETful接口，可以对接第三方接口，需要第三方鉴权
另外之后分组之后还能添加全局和局部中间件，使得功能无限扩展
如何实现：
1. 添加Group对象，访问Router
*/

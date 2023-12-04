package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
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
	router       *router
	groups       []*RouterGroup     // 保存所有groups
	htmlTemplate *template.Template // html渲染, 将所有模板加载进内存
	funcMap      template.FuncMap   // html渲染，所有自定义的模板渲染函数
}

// New函数创建一个gee.Enige
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
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

// 定义use函数，应用中间件到具体的Group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// 解析请求的路径，查找路由映射表，如果查到，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND
// day_five 当接收到具体的请求后判断适用于哪些中间件，这里简单的通过前缀来判断，得到中间件列表并赋值给c.handlers
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine // 由于在context中添加了enigne成员变量，在实例化context的时候需要给c.enigne赋值
	engine.router.handle(c)
}

// 实现模板渲染
// 解析请求的地址，映射到服务器上文件的真实地址，交给http.FileServer
// 创建静态文件的处理接口
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// 提供静态文件——Static方法暴露给用户，用来将磁盘上的文件root映射到路由上
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPath := path.Join(relativePath, "/*filepath")
	// 注册GET接口
	group.GET(urlPath, handler)
}

func (engine *Engine) SetFunMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplate = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
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

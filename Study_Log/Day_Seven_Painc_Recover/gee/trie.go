package gee

import (
	"strings"
)

/*
HTTP请求的路劲恰好是由 / 分割的多段构成
所以，每一段都可以作为前缀树的一个节点
通过树结构查询，如果中间某一层的节点都不满足条件，说明没有匹配到路由，查询结束
接下来实现的动态路由具备以下功能
- 参数匹配
- 通配 *
*/

// Trie树实现
type node struct {
	pattern  string  // 待匹配路由，例如/p/:lang
	part     string  // 路由中的一部分 :lang
	children []*node // 子节点 [doc, tutorial, intro]
	isWild   bool    // 判断是否精确匹配
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 对于路由来说，最重要的就是注册和匹配。
// 开发服务的时候，注册路由规则，映射handler
// 访问的时候匹配路由规则，查找对应的handler
// Trie树需要支持节点的插入和查询，插入——》递归查找每一层节点，没有找到就创建一个

// 注册路由，映射handler
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 递归查询每一层节点
// 匹配到了“*”或者匹配到了第len(parts)层节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		res := child.search(parts, height+1)
		if res != nil {
			return res
		}
	}
	return nil
}

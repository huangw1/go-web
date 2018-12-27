package router

import (
	"strings"
)

type node struct {
	component string
	methods map[string] Handler
	isNamedParam bool
	children []*node
}

func (n *node) addNode(method, path string, handler Handler)  {
	components := strings.Split(path[1:], "/")
	n.add(method, components, handler)
}

func (n *node) add(method string, components []string, handler Handler) {
	nd, comp, comps := n.traverse(components, nil)
	count := len(comps)
	if nd.component == comp && count == 0 {
		nd.methods[method] = handler
	} else {
		nnd := &node{
			component: comp,
			methods: make(map[string] Handler),
			isNamedParam: false,
		}
		if len(comp) > 0 && comp[0] == ':' {
			nnd.isNamedParam = true
		}
		nd.children = append(nd.children, nnd)
		if count == 0 {
			nnd.methods[method] = handler
		} else {
			nnd.add(method, comps, handler)
		}
	}
}

func (n *node) traverse(components []string, params map[string]string) (*node, string, []string) {
	component := components[0]
	if len(n.children) > 0 {
		for _, child := range n.children {
			if child.component == component || child.isNamedParam {
				if child.isNamedParam && params != nil {
					params[child.component[1:]] = component
				}
				nextComponents := components[1:]
				if len(nextComponents) > 0 {
					return child.traverse(nextComponents, params)
				}
				return child, component, nextComponents
			}
		}
	}
	return n, component, components[1:]
}
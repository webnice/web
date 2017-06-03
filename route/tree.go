package route

// Radix tree implementation below is a based on the original work by
// Armon Dadgar in https://github.com/armon/go-radix/blob/master/radix.go
// (MIT licensed).

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context"
import "gopkg.in/webnice/web.v1/method"
import (
	"net/http"
	"sort"
	"strings"
)

const (
	ntStatic   nodeTyp = iota // /string
	ntParam                   // /:variable
	ntCatchAll                // /api/v1.0/*
	//ntRegexp                  // /:id([0-9]+) or #id^[0-9]+$
)

type (
	nodeTyp        uint8
	walkFn         func(pattern string, handlers methodHandlers, subroutes Routes) bool
	methodHandlers map[method.Method]http.Handler // is a mapping of http method constants to handlers for a given route
)

// Route structure
type Route struct {
	Pattern   string
	Handlers  map[string]http.Handler
	SubRoutes Routes
}

type node struct {
	// node type
	typ nodeTyp

	// first byte of the prefix
	label byte

	// prefix is the common prefix we ignore
	prefix string

	// pattern is the computed path of prefixes
	pattern string

	// HTTP handler on the leaf node
	handlers methodHandlers

	// chi subroutes on the leaf node
	subroutes Routes

	// Child nodes should be stored in-order for iteration,
	// in groups of the node type.
	children [ntCatchAll + 1]nodes
}

type nodes []*node

// Sort the list of nodes by label
func (ns nodes) Len() int           { return len(ns) }
func (ns nodes) Less(i, j int) bool { return ns[i].label < ns[j].label }
func (ns nodes) Swap(i, j int)      { ns[i], ns[j] = ns[j], ns[i] }
func (ns nodes) Sort()              { sort.Sort(ns) }

// FindRoute Поиск роутинга и методов
func (n *node) FindRoute(ctx context.Interface, path string) methodHandlers {
	var rn *node

	// Reset the context routing pattern
	ctx.Route().Pattern("")

	// Find the routing handlers for the path
	rn = n.findRoute(ctx, path)
	if rn == nil {
		return nil
	}

	// Record the routing pattern in the request lifecycle
	if rn.pattern != "" {
		ctx.Route().Pattern(rn.pattern)
		ctx.Route().Patterns(append(ctx.Route().Patterns(), ctx.Route().Pattern()))
	}

	return rn.handlers
}

// InsertRoute Добавление роутинга и обработчика
func (n *node) InsertRoute(mtd method.Method, pattern string, handler http.Handler) *node {
	var cn, parent, child, subchild *node
	var search = pattern
	var p, commonPrefix int

	for {
		// Handle key exhaustion
		if len(search) == 0 {
			// Insert or update the node's leaf handler
			n.setHandler(mtd, handler)
			n.pattern = pattern
			return n
		}

		// Look for the edge
		parent = n
		n = n.getEdge(search[0])

		// No edge, create one
		if n == nil {
			cn = &node{label: search[0], prefix: search, pattern: pattern}
			cn.setHandler(mtd, handler)
			parent.addChild(pattern, cn)
			return cn
		}

		if n.typ > ntStatic {
			// We found a wildcard node, meaning search path starts with
			// a wild prefix. Trim off the wildcard search path and continue
			p = strings.Index(search, "/")
			if p < 0 {
				p = len(search)
			}
			search = search[p:]
			continue
		}

		// Static nodes fall below here
		// Determine longest prefix of the search key on match
		commonPrefix = n.longestPrefix(search, n.prefix)
		if commonPrefix == len(n.prefix) {
			// the common prefix is as long as the current node's prefix we're attempting to insert
			// keep the search going.
			search = search[commonPrefix:]
			continue
		}

		// Split the node
		child = &node{
			typ:    ntStatic,
			prefix: search[:commonPrefix],
		}
		parent.replaceChild(search[0], child)

		// Restore the existing node
		n.label = n.prefix[commonPrefix]
		n.prefix = n.prefix[commonPrefix:]
		child.addChild(pattern, n)

		// If the new key is a subset, add to to this node
		search = search[commonPrefix:]
		if len(search) == 0 {
			child.setHandler(mtd, handler)
			child.pattern = pattern
			return child
		}

		// Create a new edge for the node
		subchild = &node{
			typ:     ntStatic,
			label:   search[0],
			prefix:  search,
			pattern: pattern,
		}
		subchild.setHandler(mtd, handler)
		child.addChild(pattern, subchild)
		return subchild
	}
}

// Поиск шаблона
func (n *node) findPattern(pattern string) *node {
	var nn *node
	var nds nodes
	var idx int
	var xpattern string
	nn = n
	for _, nds = range nn.children {
		if len(nds) == 0 {
			continue
		}

		n = nn.getEdge(pattern[0])
		if n == nil {
			continue
		}

		idx = n.longestPrefix(pattern, n.prefix)
		xpattern = pattern[idx:]

		if len(xpattern) == 0 {
			return n
		} else if xpattern[0] == '/' && idx < len(n.prefix) {
			continue
		}

		return n.findPattern(xpattern)
	}
	return nil
}

func (n *node) isLeaf() bool { return n.handlers != nil }

func (n *node) addChild(pattern string, child *node) {
	var search string
	var i int
	var handlers methodHandlers
	var ntyp nodeTyp

	search = child.prefix

	// Find any wildcard segments
	i = strings.IndexAny(search, ":*")

	// Determine new node type
	ntyp = child.typ
	if i >= 0 {
		switch search[i] {
		case ':':
			ntyp = ntParam
		case '*':
			ntyp = ntCatchAll
		}
	}

	if i == 0 {
		// Path starts with a wildcard

		handlers = child.handlers
		child.typ = ntyp

		if ntyp == ntCatchAll {
			i = -1
		} else {
			i = strings.IndexByte(search, '/')
		}
		if i < 0 {
			i = len(search)
		}
		child.prefix = search[:i]

		if i != len(search) {
			// add edge for the remaining part, split the end
			child.handlers = nil

			search = search[i:]

			child.addChild(pattern, &node{
				typ:      ntStatic,
				label:    search[0], // this will always start with /
				prefix:   search,
				pattern:  pattern,
				handlers: handlers,
			})
		}

	} else if i > 0 {
		// Path has some wildcard

		// starts with a static segment
		handlers = child.handlers
		child.typ = ntStatic
		child.prefix = search[:i]
		child.handlers = nil

		// add the wild edge node
		search = search[i:]

		child.addChild(pattern, &node{
			typ:      ntyp,
			label:    search[0],
			prefix:   search,
			pattern:  pattern,
			handlers: handlers,
		})

	} else {
		// Path is all static
		child.typ = ntyp

	}

	n.children[child.typ] = append(n.children[child.typ], child)
	n.children[child.typ].Sort()
}

func (n *node) replaceChild(label byte, child *node) {
	const missingChild = `replacing missing child`
	var i int
	for i = 0; i < len(n.children[child.typ]); i++ {
		if n.children[child.typ][i].label == label {
			n.children[child.typ][i] = child
			n.children[child.typ][i].label = label
			return
		}
	}
	panic(missingChild)
}

func (n *node) getEdge(label byte) (ret *node) {
	var nds nodes
	var i, num int
	for _, nds = range n.children {
		num = len(nds)
		for i = 0; i < num; i++ {
			if nds[i].label == label {
				ret = nds[i]
			}
		}
	}
	return
}

func (n *node) findEdge(ntyp nodeTyp, label byte) *node {
	var nds = n.children[ntyp]
	var num = len(nds)
	var idx = int(0)
	var i, j int

	switch ntyp {
	case ntStatic:
		i, j = 0, num-1
		for i <= j {
			idx = i + (j-i)/2
			if label > nds[idx].label {
				i = idx + 1
			} else if label < nds[idx].label {
				j = idx - 1
			} else {
				i = num // breaks cond
			}
		}
		if nds[idx].label != label {
			return nil
		}
		return nds[idx]

	default:
		// TODO
		// wild nodes
		// right now we match them all, but regexp should
		// run through regexp matcher
		return nds[idx]
	}
}

// Recursive edge traversal by checking all nodeTyp groups along the way
// It's like searching through a multi-dimensional radix trie
func (n *node) findRoute(ctx context.Interface, path string) *node {
	var nn, xn, fin *node
	var nds nodes
	var ntyp nodeTyp
	var t, p int
	var label byte
	var search, xsearch string

	nn = n
	search = path
	for t, nds = range nn.children {
		ntyp = nodeTyp(t)
		if len(nds) == 0 {
			continue
		}

		// search subset of edges of the index for a matching node
		label = 0
		if search != "" {
			label = search[0]
		}

		xn = nn.findEdge(ntyp, label) // next node
		if xn == nil {
			continue
		}

		// Prepare next search path by trimming prefix from requested path
		xsearch = search
		if xn.typ > ntStatic {
			p = -1
			if xn.typ < ntCatchAll {
				p = strings.IndexByte(xsearch, '/')
			}
			if p < 0 {
				p = len(xsearch)
			}

			if xn.typ == ntCatchAll {
				ctx.Route().UrnParams().Add("*", xsearch)
			} else {
				ctx.Route().UrnParams().Add(xn.prefix[1:], xsearch[:p])
			}

			xsearch = xsearch[p:]
		} else if strings.HasPrefix(xsearch, xn.prefix) {
			xsearch = xsearch[len(xn.prefix):]
		} else {
			continue // no match
		}

		// did we find it yet?
		if len(xsearch) == 0 {
			if xn.isLeaf() {
				return xn
			}
		}

		// recursively find the next node..
		fin = xn.findRoute(ctx, xsearch)
		if fin != nil {
			// found a node, return it
			return fin
		}

		// Did not found final handler, let's remove the param here if it was set
		if xn.typ > ntStatic {
			if xn.typ == ntCatchAll {
				ctx.Route().UrnParams().Del("*")
			} else {
				ctx.Route().UrnParams().Del(xn.prefix[1:])
			}
		}
	}

	return nil
}

// finds the length of the shared prefix of two strings
func (n *node) longestPrefix(k1, k2 string) (i int) {
	var max, l int
	max = len(k1)
	if l = len(k2); l < max {
		max = l
	}
	for i = 0; i < max; i++ {
		if k1[i] != k2[i] {
			break
		}
	}
	return
}

func (n *node) setHandler(mtd method.Method, handler http.Handler) {
	var m method.Method

	if n.handlers == nil {
		n.handlers = make(methodHandlers)
	}
	if mtd.Int64()&method.Stub.Int64() == method.Stub.Int64() {
		n.handlers[method.Stub] = handler
	} else {
		n.handlers[method.Stub] = nil
	}
	if mtd&method.Any == method.Any {
		n.handlers[method.Any] = handler
		for _, m = range method.All() {
			n.handlers[m] = handler
		}
	} else {
		n.handlers[mtd] = handler
	}
}

func (n *node) isEmpty() (ret bool) {
	var nds nodes
	for _, nds = range n.children {
		if len(nds) > 0 {
			return
		}
	}
	ret = true
	return
}

func (n *node) routes() []Route {
	var rts []Route

	n.walkRoutes(n.prefix, n, func(pattern string, handlers methodHandlers, subroutes Routes) (ret bool) {
		var h http.Handler
		var mt method.Method
		var hs map[string]http.Handler
		var m string

		if handlers[method.Stub] != nil && subroutes == nil {
			return
		}
		if subroutes != nil && len(pattern) > 2 {
			pattern = pattern[:len(pattern)-2]
		}
		hs = make(map[string]http.Handler)
		if handlers[method.Any] != nil {
			hs["*"] = handlers[method.Any]
		}
		for mt, h = range handlers {
			if h == nil {
				continue
			}
			m = mt.String()
			if m == "" {
				continue
			}
			hs[m] = h
		}
		rts = append(rts, Route{pattern, hs, subroutes})
		ret = true
		return
	})

	return rts
}

func (n *node) walkRoutes(pattern string, nd *node, fn walkFn) (ret bool) {
	var nds nodes
	var nn *node
	var pat = nd.pattern

	// Visit the leaf values if any
	if (nd.handlers != nil || nd.subroutes != nil) && fn(pat, nd.handlers, nd.subroutes) {
		ret = true
		return
	}

	// Recurse on the children
	for _, nds = range nd.children {
		for _, nn = range nds {
			if n.walkRoutes(pat, nn, fn) {
				ret = true
				return
			}
		}
	}
	return
}

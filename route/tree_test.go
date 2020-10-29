package route

// Radix tree implementation below is a based on the original work by
// Armon Dadgar in https://github.com/armon/go-radix/blob/master/radix.go
// (MIT licensed).

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/webnice/web/v1/context"
	"github.com/webnice/web/v1/method"
)

var (
	emptyParams = map[string]string{}
)

func TestTree(t *testing.T) {
	hStub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hIndex := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hRobots := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hArticleList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hArticleNear := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hArticleShow := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hArticleShowRelated := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hArticleShowOpts := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hArticleSlug := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hArticleByUser := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hUserList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hUserShow := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hAdminCatchall := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hAdminAppShow := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hAdminAppShowCatchall := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hUserProfile := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hUserSuper := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hUserAll := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hHubView1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hHubView2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	hHubView3 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	//

	tr := &node{}
	_, _ = tr.InsertRoute(method.Stub, "/*", hStub)
	if routes := tr.routes(); len(routes) > 0 {
		t.Errorf("Error function routes()")
	}
	_, _ = tr.InsertRoute(method.Any, "//*", hIndex)
	subroute := New()
	subroute.Get("/users", hHubView3)
	_, _ = tr.InsertRoute(method.Get, "///*", subroute)
	if routes := tr.routes(); len(routes) != 2 || len(subroute.Routes()) > 1 {
		t.Errorf("Error function routes()")
	}

	//

	tr = &node{}
	_, _ = tr.InsertRoute(method.Get, "/aaa/*", hIndex)
	_, _ = tr.InsertRoute(method.Get, "/eee/zzz", hIndex)
	_, _ = tr.InsertRoute(method.Get, "/eee/ccc", hIndex)
	_, _ = tr.InsertRoute(method.Get, "/zzz/*", hIndex)
	_, _ = tr.InsertRoute(method.Get, "/aaa", hIndex)
	_, _ = tr.InsertRoute(method.Get, "/aaa/aaa", hIndex)
	_, _ = tr.InsertRoute(method.Any, "", hIndex)
	//

	//

	tr = &node{}
	if !tr.isEmpty() {
		t.Errorf("Error function isEmpty()")
	}
	if routes := tr.routes(); len(routes) > 0 {
		t.Errorf("Error function routes()")
	}
	_, _ = tr.InsertRoute(method.Get, "/", hIndex)
	_, _ = tr.InsertRoute(method.Get, "/robots.txt", hRobots)

	_, _ = tr.InsertRoute(method.Get, "/pages/*", hStub)

	_, _ = tr.InsertRoute(method.Get, "/articles", hArticleList)
	_, _ = tr.InsertRoute(method.Get, "/articles/", hArticleList)

	_, _ = tr.InsertRoute(method.Get, "/articles/near", hArticleNear)
	_, _ = tr.InsertRoute(method.Get, "/articles/:id", hStub)
	_, _ = tr.InsertRoute(method.Get, "/articles/:id", hArticleShow)
	_, _ = tr.InsertRoute(method.Get, "/articles/:id", hArticleShow)
	_, _ = tr.InsertRoute(method.Get, "/articles/@:user", hArticleByUser)

	_, _ = tr.InsertRoute(method.Get, "/articles/:sup/:opts", hArticleShowOpts)
	_, _ = tr.InsertRoute(method.Get, "/articles/:id/:opts", hArticleShowOpts)

	_, _ = tr.InsertRoute(method.Get, "/articles/:iffd/edit", hStub)
	_, _ = tr.InsertRoute(method.Get, "/articles/:id//related", hArticleShowRelated)
	_, _ = tr.InsertRoute(method.Get, "/articles/slug/:month/-/:day/:year", hArticleSlug)

	_, _ = tr.InsertRoute(method.Get, "/admin/user", hUserList)
	_, _ = tr.InsertRoute(method.Get, "/admin/user/", hStub) // will get replaced by next route
	_, _ = tr.InsertRoute(method.Get, "/admin/user/", hUserList)

	_, _ = tr.InsertRoute(method.Get, "/admin/user//:id", hUserShow)
	_, _ = tr.InsertRoute(method.Get, "/admin/user/:id", hUserShow)

	_, _ = tr.InsertRoute(method.Get, "/admin/applications/:id", hAdminAppShow)
	_, _ = tr.InsertRoute(method.Get, "/admin/applications/:id/*ff", hAdminAppShowCatchall)

	_, _ = tr.InsertRoute(method.Get, "/admin/*ff", hStub) // catchall segment will get replaced by next route
	_, _ = tr.InsertRoute(method.Get, "/admin/*", hAdminCatchall)

	_, _ = tr.InsertRoute(method.Get, "/users/:userID/profile", hUserProfile)
	_, _ = tr.InsertRoute(method.Get, "/users/master/*", hUserSuper)
	_, _ = tr.InsertRoute(method.Get, "/users/*", hUserAll)

	_, _ = tr.InsertRoute(method.Get, "/hubs/:hubID/view", hHubView1)
	_, _ = tr.InsertRoute(method.Get, "/hubs/:hubID/view/*", hHubView2)

	subroute = New()
	subroute.Get("/users", hHubView3)
	_, _ = tr.InsertRoute(method.Get, "/hubs/:hubID/*", subroute)
	_, _ = tr.InsertRoute(method.Get, "/hubs/:hubID/users", hHubView3)

	if tr.isEmpty() {
		t.Errorf("Error function isEmpty()")
	}

	tests := []struct {
		r string            // input request path
		h http.Handler      // output matched handler
		p map[string]string // output params
	}{
		{r: "/", h: hIndex, p: emptyParams},
		{r: "/robots.txt", h: hRobots, p: emptyParams},

		{r: "/pages", h: nil, p: emptyParams},
		{r: "/pages/", h: hStub, p: map[string]string{"*": ""}},
		{r: "/pages/yes", h: hStub, p: map[string]string{"*": "yes"}},

		{r: "/articles", h: hArticleList, p: emptyParams},
		{r: "/articles/", h: hArticleList, p: emptyParams},
		{r: "/articles/near", h: hArticleNear, p: emptyParams},
		{r: "/articles/neard", h: hArticleShow, p: map[string]string{"id": "neard"}},
		{r: "/articles/123", h: hArticleShow, p: map[string]string{"id": "123"}},
		{r: "/articles/123/456", h: hArticleShowOpts, p: map[string]string{"id": "123", "opts": "456"}},
		{r: "/articles/@peter", h: hArticleByUser, p: map[string]string{"user": "peter"}},
		{r: "/articles/22//related", h: hArticleShowRelated, p: map[string]string{"id": "22"}},
		{r: "/articles/111/edit", h: hStub, p: map[string]string{"id": "111"}},
		{r: "/articles/slug/sept/-/4/2015", h: hArticleSlug, p: map[string]string{"month": "sept", "day": "4", "year": "2015"}},
		{r: "/articles/:id", h: hArticleShow, p: map[string]string{"id": ":id"}}, // TODO review goji?

		{r: "/admin/user", h: hUserList, p: emptyParams},
		{r: "/admin/user/", h: hUserList, p: emptyParams},
		{r: "/admin/user/1", h: hUserShow, p: map[string]string{"id": "1"}}, // hmmm.... TODO, review
		{r: "/admin/user//1", h: hUserShow, p: map[string]string{"id": "1"}},
		{r: "/admin/hi", h: hAdminCatchall, p: map[string]string{"*": "hi"}},
		{r: "/admin/lots/of/:fun", h: hAdminCatchall, p: map[string]string{"*": "lots/of/:fun"}},
		{r: "/admin/applications/333", h: hAdminAppShow, p: map[string]string{"id": "333"}},
		{r: "/admin/applications/333/woot", h: hAdminAppShowCatchall, p: map[string]string{"id": "333", "*": "woot"}},

		{r: "/hubs/123/view", h: hHubView1, p: map[string]string{"hubID": "123"}},
		{r: "/hubs/123/view/index.html", h: hHubView2, p: map[string]string{"hubID": "123", "*": "index.html"}},
		{r: "/hubs/123/users", h: hHubView3, p: map[string]string{"hubID": "123"}},

		{r: "/users/123/profile", h: hUserProfile, p: map[string]string{"userID": "123"}},
		{r: "/users/master/123/okay/yes", h: hUserSuper, p: map[string]string{"*": "123/okay/yes"}},
		{r: "/users/123/okay/yes", h: hUserAll, p: map[string]string{"*": "123/okay/yes"}},
	}

	for i, tt := range tests {
		rctx := context.New()
		handlers := tr.FindRoute(rctx, tt.r) //, params)
		handler := handlers[method.Get]
		params := rctx.Route().Params()
		if fmt.Sprintf("%v", tt.h) != fmt.Sprintf("%v", handler) {
			t.Errorf("Input [%d]: find '%s' expecting handler:%v , got:%v", i, tt.r, tt.h, handler)
		}
		for key := range tt.p {
			if !params.Has(key) {
				t.Errorf("Input [%d]: find '%s' key not found", i, tt.r)
				break
			}
		}
	}
	if routes := tr.routes(); len(routes) == 0 {
		t.Errorf("Error function routes()")
	}
}

//func TestTreeErrors(t *testing.T) {
//	//hStub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
//	hIndex := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
//
//	tr := &node{}
//	tr.InsertRoute(method.Get, "/", hIndex)
//
//	subroute := New()
//	subroute.Get("/pages", hIndex)
//	tr.InsertRoute(method.Get, "/*", subroute)
//
//	tr.InsertRoute(method.Get, "/pages/*", hIndex)
//
//	tests := []struct {
//		r string            // input request path
//		h http.Handler      // output matched handler
//		p map[string]string // output params
//	}{
//		{r: "/", h: hIndex, p: emptyParams},
//
//		//		{r: "/pages", h: nil, p: emptyParams},
//		//		{r: "/pages/", h: hStub, p: map[string]string{"*": ""}},
//		//		{r: "/pages/yes", h: hStub, p: map[string]string{"*": "yes"}},
//	}
//
//	for i, tt := range tests {
//		rctx := context.New()
//		handlers := tr.FindRoute(rctx, tt.r) //, params)
//		handler, _ := handlers[method.Get]
//		params := rctx.Route().Params()
//		if fmt.Sprintf("%v", tt.h) != fmt.Sprintf("%v", handler) {
//			t.Errorf("Input [%d]: find '%s' expecting handler:%v , got:%v", i, tt.r, tt.h, handler)
//		}
//		for key := range tt.p {
//			if !params.Has(key) {
//				t.Errorf("Input [%d]: find '%s' key not found", i, tt.r)
//				break
//			}
//		}
//	}
//
//}

//func debugPrintTree(parent int, i int, n *node, label byte) bool {
//	numEdges := 0
//	for _, nds := range n.children {
//		numEdges += len(nds)
//	}
//
//	if n.handlers != nil {
//		log.Printf("[node %d parent:%d] typ:%d prefix:%s label:%s numEdges:%d isLeaf:%v handler:%v\n", i, parent, n.typ, n.prefix, string(label), numEdges, n.isLeaf(), n.handlers)
//	} else {
//		log.Printf("[node %d parent:%d] typ:%d prefix:%s label:%s numEdges:%d isLeaf:%v\n", i, parent, n.typ, n.prefix, string(label), numEdges, n.isLeaf())
//	}
//
//	parent = i
//	for _, nds := range n.children {
//		for _, e := range nds {
//			i++
//			if debugPrintTree(parent, i, e, e.label) {
//				return true
//			}
//		}
//	}
//	return false
//}

//func BenchmarkTreeGet(b *testing.B) {
//	h1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
//	h2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
//
//	tr := &node{}
//	tr.InsertRoute(mGET, "/", h1)
//	tr.InsertRoute(mGET, "/ping", h2)
//	tr.InsertRoute(mGET, "/pingall", h2)
//	tr.InsertRoute(mGET, "/ping/:id", h2)
//	tr.InsertRoute(mGET, "/ping/:id/woop", h2)
//	tr.InsertRoute(mGET, "/ping/:id/:opt", h2)
//	tr.InsertRoute(mGET, "/pings", h2)
//	tr.InsertRoute(mGET, "/hello", h1)
//
//	b.ReportAllocs()
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		// params := map[string]string{}
//		mctx := NewRouteContext()
//		tr.FindRoute(mctx, "/ping/123/456")
//		// tr.Find("/pings", params)
//	}
//}

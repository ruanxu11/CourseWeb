package routerKL

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
)

// 路由定义
type route struct {
	pattern     string // 正则表达式
	method      string
	httpHandler func(w http.ResponseWriter, r *http.Request)
}

type routerKL struct {
	routes        []route
	listeningPort net.Listener
}

var showLog = false

// 使用正则路由转发
func (r *routerKL) route(w http.ResponseWriter, req *http.Request) {
	requestPath := req.URL.Path
	req.ParseForm()
	if showLog {
		for k, v := range req.Form {
			fmt.Print(k, ": ")
			fmt.Println(strings.Join(v, ", "))
		}
		log.Println(req.Method + " " + requestPath)
	}
	isFound := false
	for i := 0; i < len(r.routes); i++ {
		route := r.routes[i]
		if route.method == "Handle" || route.method == req.Method {
			reg, err := regexp.Compile("^" + route.pattern + "$")
			if err != nil {
				log.Println(err)
			}
			if reg.MatchString(requestPath) {
				isFound = true
				route.httpHandler(w, req)
				break
			}
		}
	}

	if !isFound {
		// 未匹配到路由
		fmt.Fprint(w, "404 Page Not Found!")
	}
}

func Handle(pattern string, handler interface{}) {
	r.addRoute(pattern, "Handle", handler)
}

func Get(pattern string, handler interface{}) {
	r.addRoute(pattern, "GET", handler)
}

func Post(pattern string, handler interface{}) {
	r.addRoute(pattern, "Post", handler)
}

func (r *routerKL) addRoute(pattern string, method string, handler interface{}) {
	r.routes = append(r.routes, route{pattern: pattern, method: method, httpHandler: handler.(func(w http.ResponseWriter, r *http.Request))})
}

var r = routerKL{}

func Run(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		r.route(w, req)
	})
	log.Println("Listening on " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func RunWithLog(addr string) {
	showLog = true
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		r.route(w, req)
	})
	log.Println("Listening on " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func ParamGet(req *http.Request, name string) string {
	query := req.URL.Query()
	if params, ok := query[name]; ok {
		return params[0]
	}
	return ""
}

func ParamGetX(req *http.Request, name string, x int) string {
	query := req.URL.Query()
	if params, ok := query[name]; ok && len(params) > x {
		return params[x]
	}
	return ""
}

func ParamPost(req *http.Request, name string) string {
	if params, ok := req.PostForm[name]; ok {
		return params[0]
	}
	return ""
}

func ParamPostX(req *http.Request, name string, x int) string {
	if params, ok := req.PostForm[name]; ok && len(params) > x {
		return params[x]
	}
	return ""
}

func Param(req *http.Request, name string) string {
	if params, ok := req.Form[name]; ok {
		return params[0]
	}
	return ""
}

func ParamX(req *http.Request, name string, x int) string {
	if params, ok := req.Form[name]; ok && len(params) > x {
		return params[x]
	}
	return ""
}

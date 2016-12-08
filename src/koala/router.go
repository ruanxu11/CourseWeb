package koala

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

// 路由定义
type Route struct {
	slice   []string
	method  string
	handler func(p *Params, w http.ResponseWriter, r *http.Request)
}

type Params struct {
	Param     map[string][]string
	ParamGet  map[string][]string
	ParamPost map[string][]string
	ParamUrl  map[string]string
}

type App struct {
	routes        []Route
	listeningPort net.Listener
}

var showLog = false

// 使用正则路由转发
func (app *App) route(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	p := Params{
		ParamGet:  r.URL.Query(),
		ParamPost: r.PostForm,
		Param:     r.Form,
		ParamUrl:  make(map[string]string),
	}

	requestPath := r.URL.Path
	r.ParseForm()
	if showLog {
		log.Println(r.Method + " " + requestPath)
	}
	isFound := false
	for i := 0; i < len(app.routes); i++ {
		route := app.routes[i]
		if route.method == "Handle" || route.method == r.Method {
			url := strings.Split(requestPath, "/")[1:]
			if len(url) != len(route.slice) {
				continue
			}
			matched := true
			for i := 0; i < len(route.slice); i++ {
				if route.slice[i][0] == ':' {
					p.ParamUrl[route.slice[i][1:]] = url[i]
				} else if route.slice[i] == url[i] {
					continue
				} else {
					matched = false
					break
				}
			}
			if !matched {
				continue
			}

			if showLog {
				log.Println(route.slice)
				fmt.Print("get: ")
				fmt.Println(p.ParamGet)
				fmt.Print("post: ")
				fmt.Println(p.ParamPost)
				fmt.Print("url: ")
				fmt.Println(p.ParamUrl)
			}
			isFound = true
			route.handler(&p, w, r)
			break

		}
	}

	if !isFound {
		// 未匹配到路由
		fmt.Fprint(w, "404 Page Not Found!")
	}
}

func Handle(pattern string, handler interface{}) {
	app.addRoute(pattern, "Handle", handler)
}

func Get(pattern string, handler interface{}) {
	app.addRoute(pattern, "GET", handler)
}

func Post(pattern string, handler interface{}) {
	app.addRoute(pattern, "Post", handler)
}

func (app *App) addRoute(pattern string, method string, handler interface{}) {
	slice := strings.Split(pattern, "/")
	app.routes = append(app.routes, Route{slice: slice[1:], method: method, handler: handler.(func(p *Params, w http.ResponseWriter, r *http.Request))})
}

var app = App{}

func Run(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.route(w, r)
	})
	log.Println("Listening on " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func RunWithLog(addr string) {
	showLog = true
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.route(w, r)
	})
	log.Println("Listening on " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

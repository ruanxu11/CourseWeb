package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"labix.org/v2/mgo/bson"

	"github.com/gorilla/sessions"
)

func renderHTML(w http.ResponseWriter, file string, data interface{}) {
	t, err := template.New(file).ParseFiles("static/views/" + file)
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, data)
}

var sessionStore = sessions.NewCookieStore([]byte("3140102431"))

func loginCheck(collection string, sid string, password string) (bool, string) {
	id, _ := strconv.Atoi(sid)
	person, err := mgoFind(collection, bson.M{
		"_id":      id,
		"password": password,
	}, 0, 0)
	log.Println(person)
	if err != nil {
		return false, ""
	}
	return true, person["name"].(string)
}

func routerInit() {
	// files
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("static/images"))))
	http.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("static/views"))))
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("static/template"))))

	// index
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fileData, _ := ioutil.ReadFile("./readme.md")
		readme := strings.Replace(string(fileData), "\n", "<br>", -1)
		readme = strings.Replace(readme, "\t", "    ", -1)
		readme = strings.Replace(readme, " ", "&nbsp", -1)
		args := map[string]template.HTML{
			"title":  template.HTML("浙江大学课程网站系统"),
			"readme": template.HTML(string(readme)),
		}
		renderHTML(w, "index.html", args)
	})

	// /login
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Println("method:", r.Method) //获取请求的方法
		ID := r.Form["ID"][0]
		Password := r.Form["Password"][0]
		Type := r.Form["Type"][0]
		fmt.Println("ID", r.Form["ID"])
		fmt.Println("Password", r.Form["Password"])
		for k, v := range r.Form {
			fmt.Print("key:", k, "; ")
			fmt.Println("val:", strings.Join(v, ""))
		}
		// if r.Method == "GET" {
		// 	query := r.URL.Query()
		// 	fmt.Println("ID", query["ID"])
		// 	fmt.Println("Password", query["Password"])
		// 	fmt.Println("Type", query["Type"])
		// } else if r.Method == "POST" {
		// }
		vaild, Name := loginCheck(Type, ID, Password)
		log.Println(Name)
		if vaild {
			session, err := sessionStore.Get(r, "sessionID")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if session.IsNew {
				log.Println("登陆成功")
				session.Values["ID"] = ID
				session.Values["Name"] = Name
				session.Values["Password"] = Password
				session.Values["Type"] = Type
				session.Save(r, w)
				w.Write([]byte("登陆成功\n欢迎," + Name + "同学\n"))
			} else {
				w.Write([]byte("您已经登陆了," + Name + "同学\n"))
			}
		} else {
			w.Write([]byte("账号或密码错误"))
		}
	})
}

// searchResults, _ := mgoFindAll("students", nil, 0, 0)
// json, _ := json.Marshal(searchResults)
// w.Write([]byte(json))
// var err error
// if err != nil {
// 	json, _ := json.Marshal(err)
// 	w.Write([]byte("err\n"))
// 	w.Write([]byte(json))
// 	log.Println(err)
// } else {
// 	w.Write([]byte("登陆成功"))
// 	w.Write([]byte("您已经登陆了," + Name + "同学"))
// }

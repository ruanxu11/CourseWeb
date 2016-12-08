package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"koala"
)

func renderHTML(w http.ResponseWriter, file string, data interface{}) {
	t, err := template.New(file).ParseFiles("static/views/" + file)
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, data)
}

func RouteStart() {
	// files
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("static/images"))))
	http.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("static/views"))))
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("static/template"))))

	accountHandler()
	classHandlers()

	koala.Handle("/", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
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

	koala.RunWithLog(":2333")
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

package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func renderHTML(w http.ResponseWriter, file string, data interface{}) {
	t, err := template.New(file).ParseFiles("static/views/" + file)
	checkErr(err)
	t.Execute(w, data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fileData, _ := ioutil.ReadFile("./readme.md")
	readme := strings.Replace(string(fileData), "\n", "<br>", -1)
	readme = strings.Replace(readme, "\t", "    ", -1)
	readme = strings.Replace(readme, " ", "&nbsp", -1)
	args := map[string]template.HTML{
		"title":  template.HTML("浙江大学课程网站系统"),
		"readme": template.HTML(string(readme)),
	}
	renderHTML(w, "index.html", args)
}

func routerInit() {
	// files
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("static/images"))))
	http.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("static/views"))))
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("static/template"))))

	// website
	http.HandleFunc("/", indexHandler)
}

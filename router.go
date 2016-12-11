package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"koala"
	"log"
)

func RouteStart() {
	// files
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("static/images"))))
	http.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("static/views"))))
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("static/template"))))
	http.Handle("/material/", http.StripPrefix("/material/", http.FileServer(http.Dir("static/upload/material"))))

	userHandlers()
	classHandlers()
	courseHandlers()

	koala.Get("/forget/password/id", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		koala.Render(w, "forgetPassword.html", map[string]interface{}{
			"title": courseWeb,
		})
	})

	koala.Post("/forget/password/questions", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		collection := p.Param["collection"][0]
		id := p.Param["id"][0]
		securityQuestions, err := getSecurityQuestions(collection, id)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/forget/password/id", "不存在该账号", "error")
		} else {
			koala.Render(w, "forgetPasswordPost.html", map[string]interface{}{
				"title":      courseWeb,
				"collection": collection,
				"id":         id,
				"question1":  securityQuestions[0],
				"question2":  securityQuestions[1],
				"question3":  securityQuestions[2],
			})
		}
	})

	koala.Get("/", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		courses, err := getDistinctCourse()
		if err != nil {
			log.Println(err)
		}
		fileData, _ := ioutil.ReadFile("./readme.md")
		readme := strings.Replace(string(fileData), "\n", "<br>", -1)
		readme = strings.Replace(readme, "\t", "    ", -1)
		readme = strings.Replace(readme, " ", "&nbsp", -1)
		koala.Render(w, "index.html", map[string]interface{}{
			"title":   courseWeb,
			"courses": courses,
			"readme":  template.HTML(string(readme)),
		})
	})
	koala.RunWithLog(":2333")
}

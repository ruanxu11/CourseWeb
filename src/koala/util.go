package koala

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var base64coder = base64.StdEncoding

func HashString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func Render(w http.ResponseWriter, file string, data interface{}) {
	t, err := template.New(file).ParseFiles("static/views/" + file)
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, data)
}

func RelocationSweet(w http.ResponseWriter, url string, title string, Type string) {
	Render(w, "relocation.html", map[string]interface{}{
		"title": title,
		"Type":  Type,
		"url":   url,
	})
}

func Relocation(w http.ResponseWriter, url string, title string, Type string) {
	t, err := template.New("x").Parse("<script>alert('" + title + "');window.location.href='" + url + "';</script>")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

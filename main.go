package main

import (
	"log"
	"net/http"
)

func main() {
	// addStudents()
	routerInit()
	go HttpGet()
	go HttpPost()
	if err := http.ListenAndServe(":2333", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

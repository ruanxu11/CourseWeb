package main

import (
	"log"
	"net/http"
)

func main() {
	routerInit()

	// http
	if err := http.ListenAndServe(":2333", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	// if err := http.ListenAndServeTLS(":8080", "server.pem", "server.key", nil); err != nil {
	// 	log.Fatal("ListenAndServe:", err)
	// }
}

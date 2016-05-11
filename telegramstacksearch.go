package main

import (
	//"github.com/Ahineya/stacksearch/api"
	"net/http"
	"fmt"
	"log"
)

func main() {
	http.HandleFunc("/", index)

	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.URL.Path == "/favicon.ico" {
		return
	}

	fmt.Println(r.Form)
	fmt.Println(r.URL.Path)

	fmt.Fprintf(w, "It works!")
}
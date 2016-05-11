package main

import (
	"github.com/Ahineya/telegramstacksearch/api"
	"net/http"
	"fmt"
	"log"
	"os"
)

func main() {
	http.HandleFunc("/", index)

	port := os.Getenv("PORT")

	if port == "" {
		port = "9090"
	}

	fmt.Println("Using port: ", port)

	err := http.ListenAndServe(":" + port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		os.Exit(1)
	}

}

func index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.URL.Path == "/favicon.ico" {
		return
	}

	fmt.Println(r.Form)
	fmt.Println(r.URL.Path)

	if len(r.Form["query"]) == 0 {
		fmt.Fprintf(w, "It works! Specify the query GET parameter")
	} else {
		response, err := api.GetAnswer(r.Form["query"][0])
		if err != nil {
			fmt.Fprintf(w, "something is broken")
		} else {
			fmt.Fprintf(w, response)
		}
	}

}
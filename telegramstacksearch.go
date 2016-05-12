package main

import (
	//"github.com/Ahineya/telegramstacksearch/api"
	"github.com/Ahineya/telegramstacksearch/telegramapi"
	"net/http"
	"fmt"
	//"log"
	"os"
	"log"
	"encoding/json"
	//"github.com/Ahineya/telegramstacksearch/api"
	"strings"
)

func main() {
	http.HandleFunc("/", index)

	port := os.Getenv("PORT")

	if port == "" {
		port = "9090"
	}

	fmt.Println("Using port: ", port)

	/*messages, err := telegramapi.GetMessages()
	if err != nil {
		log.Fatal("TelegramAPI: ", err)
		os.Exit(1)
	}*/

	//telegramapi.SendMessage(messages.Result[len(messages.Result) - 1].Message.Chat.Id, "Hello from GO")
	err := telegramapi.SetHook()
	if err != nil {
		log.Fatal("Setting Hook: ", err)
		os.Exit(1)
	}

	err = http.ListenAndServe(":" + port, nil)

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

	decoder := json.NewDecoder(r.Body)

	var t telegramapi.Update
	err := decoder.Decode(&t)

	if err != nil {
		panic(err)
	}

	//fmt.Println(t.Message.Chat.Id)

	if string(t.Message.Text[0]) == "/" {
		tokens := strings.Fields(t.Message.Text)
		command := tokens[0]
		args := tokens[1:]

		switch command {
		case "/start":
			telegramapi.SendMessage(t.Message.Chat.Id, "Welcome to Stackoverflow Search Bot!")
		case "/help":
			telegramapi.SendMessage(t.Message.Chat.Id, "Help: Stackoverflow Search Bot!")
		default:
			telegramapi.SendMessage(t.Message.Chat.Id, command + ":" + strings.Join(args, ","))
		}

	}

	telegramapi.SendMessage(t.Message.Chat.Id, "Hello from GO")

	/*if len(r.Form["query"]) == 0 {
		fmt.Fprintf(w, "It works! Specify the query GET parameter")
	} else {
		response, err := api.GetAnswer(r.Form["query"][0])
		if err != nil {
			fmt.Fprintf(w, "something is broken")
		} else {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, response)
		}
	}*/

}
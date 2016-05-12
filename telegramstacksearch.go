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
	"github.com/Ahineya/telegramstacksearch/api"
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

	decoder := json.NewDecoder(r.Body)

	var t telegramapi.Update
	err := decoder.Decode(&t)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if string(t.Message.Text[0]) == "/" {
		tokens := strings.Fields(t.Message.Text)
		command := tokens[0]
		args := tokens[1:]

		switch command {
		case "/start":
			telegramapi.SendMessage(t.Message.Chat.Id, "Welcome to Stackoverflow Search Bot!")
		case "/help":
			telegramapi.SendMessage(t.Message.Chat.Id, "Help: Stackoverflow Search Bot!")
		case "/lucky":
			if len(args) > 0 {
				response, err := api.GetAnswer(strings.Join(args, "%20"))
				if err != nil {
					telegramapi.SendMessage(t.Message.Chat.Id, "Got an error: " + err.Error())
				} else {
					telegramapi.PostMessage(response, t.Message.Chat.Id)
				}
			} else {
				telegramapi.SendMessage(t.Message.Chat.Id, "Please, specify the search query")
			}
		default:
			telegramapi.SendMessage(t.Message.Chat.Id, command + ":" + strings.Join(args, ","))
		}

	}

	w.Write([]byte {42})

}
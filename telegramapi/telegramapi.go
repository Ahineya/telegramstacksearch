package telegramapi

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"net/url"
	"bytes"
	"regexp"
)

type TelegramMessage struct {
	Text string `json:"text"`
	Chat struct {
		Id int `json:"id"`
	     } `json:"chat"`
}

type OutgoingTelegramMessage struct {
	Text string `json:"text"`
	ChatId int `json:"chat_id"`
	ParseMode string `json:"parse_mode"`
}

type Update struct {
	Id int `json:"update_id"`
	Message TelegramMessage `json:"message"`
}

type ApiResult struct {
	Ok bool `json:"ok"`
	Result []Update `json:"result"`
}

func SetHook() error {
	token := os.Getenv("BOT_TOKEN")
	api_url := "https://api.telegram.org/bot" + token + "/"

	_, err := http.Get(api_url + "setWebhook?url=https://telegramstacksearch.herokuapp.com/")
	return err
}

func GetMessages() (*ApiResult, error) {
	token := os.Getenv("BOT_TOKEN")
	api_url := "https://api.telegram.org/bot" + token + "/"

	resp, err := http.Get(api_url + "getUpdates")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var messages ApiResult;
	err = json.Unmarshal(contents, &messages)

	fmt.Printf("%s", messages)

	return &messages, err
}

func SendMessage(chatId int, text string) (bool, error) {
	token := os.Getenv("BOT_TOKEN")
	api_url := "https://api.telegram.org/bot" + token + "/"

	encodedText, err := url.Parse(text)
	if err != nil {
		return false, err
	}

	fmt.Printf("%d:%s\n", chatId, encodedText.String())

	resp, err := http.Get(api_url + "sendMessage?chat_id=" + strconv.Itoa(chatId) + "&text=" + encodedText.String())

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var messages ApiResult;
	err = json.Unmarshal(contents, &messages)
	if err != nil {
		return false, err
	}

	return messages.Ok, nil
}

func PostMessage(response string, chatId int) {
	if (len(response) > 4000) {
		response = response[0:4000]
	}

	r := regexp.MustCompile("<p>")
	fmt.Println(r)
	response = r.ReplaceAllString(response, "")

	r = regexp.MustCompile("</p>")
	fmt.Println(r)
	response = r.ReplaceAllString(response, "")


	a := OutgoingTelegramMessage{Text: response, ChatId: chatId, ParseMode: "HTML"}

	b, err := json.Marshal(&a)
	if err != nil {
		fmt.Println(err)
		return
	}

	token := os.Getenv("BOT_TOKEN")
	api_url := "https://api.telegram.org/bot" + token + "/sendMessage"

	req, err := http.NewRequest("POST", api_url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))



	//fmt.Println(string(b))
}

/*
r = requests.get(URL + "?offset=%s" % (last + 1))
	if r.status_code == 200:
	for message in r.json()["result"]:
	last = int(message["update_id"])
	requests.post("http://localhost:8888/",
		data=json.dumps(message),
	headers={'Content-type': 'application/json',
	'Accept': 'text/plain'}
)
 */
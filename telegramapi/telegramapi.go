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
	"github.com/kennygrant/sanitize"
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
	Description string `json:"description"`
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

func PostMessage(response string, chatId int, mode string) {
	if (len(response) > 3500) {
		response = response[0:3500]
	}

	response, err := sanitize.HTMLAllowing(response, []string{"b", "strong", "i", "em", "a", "code"})
	response = replace(response, "<strong><em>", "<strong>")
	response = replace(response, "</em></strong>", "</strong>")

	fmt.Println(response)

	a := OutgoingTelegramMessage{Text: response, ChatId: chatId, ParseMode: mode}

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

	var result ApiResult;
	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}

	if !result.Ok {
		PostMessage("Some troubles with parsing answer HTML. Try another request till this will be fixed. \n\n\n Debug info: \n" + response + "\n" + result.Description + "\n", chatId, "")
	}

}

func strip(str string, tag string) string {
	r := regexp.MustCompile(tag)
	fmt.Println(r)
	return r.ReplaceAllString(str, "")
}

func replace(str string, tag string, replacer string) string {
	r := regexp.MustCompile(tag)
	fmt.Println(r)
	return r.ReplaceAllString(str, replacer)
}

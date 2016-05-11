package api

import (
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"net/url"
	"errors"
)

const (
	get_questions_url = "http://api.stackexchange.com/2.2/search?order=desc&sort=relevance&intitle={query}&site=stackoverflow"
	get_answers_url = "https://api.stackexchange.com/2.2/answers/{answer_id}?site=stackoverflow&filter=withbody"
)

type StackQuestionItem struct {
	AcceptedAnswerID int `json:"accepted_answer_id"`
	Title string `json:"title"`
}

type StackAnswerItem struct {
	Body string `json:"body"`
}

type StackQuestionResponse struct {
	Items []StackQuestionItem `json:"items"`
}

type StackAnswerResponse struct {
	Items []StackAnswerItem `json:"items"`
	ErrorName string `json:"error_name"`
}

func GetAnswer(query string) (string, error) {
	encodedQuery, err := url.Parse(query)
	if err != nil {
		return "", err
	}

	stackQuestionResponse, err := getStackOverflowQuestion(encodedQuery.String())
	if err != nil {
		return "", err
	}

	if len(stackQuestionResponse.Items) == 0 {
		return "", errors.New("Results not found for query: " + query)
	}

	answerId := strconv.Itoa(stackQuestionResponse.Items[0].AcceptedAnswerID)
	stackAnswerResponse, err := getStackOverflowAnswer(answerId)
	if err != nil {
		return "", err
	}

	if len(stackAnswerResponse.Items) == 0 {
		return "", errors.New("Results not found for query: " + query)
	}

	answer := stackQuestionResponse.Items[0].Title + "\n\n" + string(stackAnswerResponse.Items[0].Body)

	return answer, err
}

func getStackOverflowQuestion(query string) (*StackQuestionResponse, error) {
	resp, err := http.Get(strings.Replace(get_questions_url, "{query}", query, 1))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var stackQuestionResponse StackQuestionResponse
	err = json.Unmarshal(contents, &stackQuestionResponse)

	return &stackQuestionResponse, err
}

func getStackOverflowAnswer(answerId string) (*StackAnswerResponse, error) {
	resp, err := http.Get(strings.Replace(get_answers_url, "{answer_id}", answerId, 1))
	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var stackAnswerResponse StackAnswerResponse
	err = json.Unmarshal(contents, &stackAnswerResponse)

	return &stackAnswerResponse, err
}

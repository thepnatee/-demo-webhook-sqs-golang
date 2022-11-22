package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type LineError struct {
	Message string `json:"message"`
	Details []struct {
		Message  string `json:"message"`
		Property string `json:"property"`
	} `json:"details"`
}

func Reply(replyToken string, messages map[string]interface{}) (bool, *LineError) {

	url := "https://api.line.me/v2/bot/message/reply"
	method := "POST"

	fmt.Println("messages", messages)
	bodyMessage := make(map[string]interface{})
	bodyMessage["replyToken"] = replyToken
	bodyMessage["messages"] = messages["messages"]

	dataString, _ := json.Marshal(bodyMessage)

	var jsonData = []byte(string(dataString))

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	req.Header.Add("Authorization", "Bearer "+os.Getenv("CHANNEL_TOKEN"))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return false, nil
	}
	defer res.Body.Close()

	responseBody, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 400 {
		lineErr := new(LineError)
		json.Unmarshal(responseBody, lineErr)
		if lineErr.Details != nil {
			return false, lineErr
		} else {
			return true, nil
		}

	}
	return true, nil
}
func ValidateReply(message map[string]interface{}) (bool, *LineError) {

	url := "https://api.line.me/v2/bot/message/validate/reply"
	method := "POST"
	dataString, _ := json.Marshal(message)
	var jsonData = []byte(string(dataString))
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	req.Header.Add("Authorization", "Bearer "+os.Getenv("CHANNEL_TOKEN"))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return false, nil
	}
	defer res.Body.Close()

	responseBody, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 400 {
		lineErr := new(LineError)

		json.Unmarshal(responseBody, lineErr)
		if lineErr.Details != nil {
			return false, lineErr
		} else {
			return true, nil
		}

	}
	return true, nil
}

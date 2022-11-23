package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func ImportFileJson1(filename string) map[string]interface{} {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()

	// convert json to interface
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var messages map[string]interface{}
	json.Unmarshal([]byte(byteValue), &messages)
	return messages
}

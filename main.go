package main

import (
	"demo-sqs/util"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type LineMessage struct {
	Destination string `json:"destination"`
	Events      []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Mode       string `json:"mode"`
		Timestamp  int64  `json:"timestamp"`
		Source     struct {
			Type   string `json:"type"`
			UserID string `json:"userId"`
		} `json:"source"`
		WebhookEventID  string `json:"webhookEventId"`
		DeliveryContext struct {
			IsRedelivery bool `json:"isRedelivery"`
		} `json:"deliveryContext"`
		Message struct {
			ID     string `json:"id"`
			Type   string `json:"type"`
			Text   string `json:"text"`
			Emojis []struct {
				Index     int    `json:"index"`
				Length    int    `json:"length"`
				ProductID string `json:"productId"`
				EmojiID   string `json:"emojiId"`
			} `json:"emojis"`
			Mention struct {
				Mentionees []struct {
					Index  int    `json:"index"`
					Length int    `json:"length"`
					UserID string `json:"userId"`
				} `json:"mentionees"`
			} `json:"mention"`
		} `json:"message"`
	} `json:"events"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()

	go util.SQSConnect()
	go util.Read()

	app.Post("/webhook", func(c *fiber.Ctx) error {

		events := new(LineMessage)
		err = c.BodyParser(events)
		if err != nil {
			log.Println(err)
		}

		for _, event := range events.Events {
			queue := make(map[string]interface{})
			queue["replyToken"] = event.ReplyToken
			if event.Message.Type == "text" {

				if event.Message.Text == "text1" {
					// start import file json object
					jsonFile, err := os.Open("message_demo_1.json")
					if err != nil {
						log.Println(err)
					}
					defer jsonFile.Close()

					// convert json to interface
					byteValue, _ := ioutil.ReadAll(jsonFile)
					var messages map[string]interface{}
					json.Unmarshal([]byte(byteValue), &messages)

					// LINE API Validate Object
					chkValid, errMessage := util.ValidateReply(messages)

					// condition true
					if chkValid {
						queue["messages"] = messages
						info, _ := json.Marshal(queue)

						// deliver message to queue
						util.SQSWriter(string(info))
					} else {
						log.Println(*errMessage)
					}

				}
				if event.Message.Text == "text2" {

					// start import file json error object
					jsonFile, err := os.Open("message_demo_2.json")
					if err != nil {
						log.Println(err)
					}
					defer jsonFile.Close()

					// convert json to interface
					byteValue, _ := ioutil.ReadAll(jsonFile)
					var messages map[string]interface{}
					json.Unmarshal([]byte(byteValue), &messages)

					// LINE API Validate Object
					chkValid, errMessage := util.ValidateReply(messages)

					// condition false
					if chkValid {
						queue["messages"] = messages
						info, _ := json.Marshal(queue)
						// deliver message to queue
						util.SQSWriter(string(info))
					} else {
						log.Println(*errMessage)
					}

				}
			}

		}

		return c.SendStatus(200)
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.SendString("Not found.")
	})

	appPort := os.Getenv("PORT")
	log.Println("Application is starting at port :", appPort)
	app.Listen(":" + os.Getenv("PORT"))
}

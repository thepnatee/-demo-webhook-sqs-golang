package util

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var sqsSvc *sqs.SQS

func SQSConnect() *sqs.SQS {

	accessKeyId := os.Getenv("CUSTOM_AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("CUSTOM_AWS_SECRET_ACCESS_KEY")

	sess := session.New(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_S3_REGION")),
		Credentials: credentials.NewStaticCredentials(accessKeyId, secretAccessKey, ""),
		MaxRetries:  aws.Int(5),
	})

	sqsSvc = sqs.New(sess)

	return sqsSvc
}

func SQSWriter(str string) {
	send_params := &sqs.SendMessageInput{
		MessageBody:  aws.String(str),
		QueueUrl:     aws.String(os.Getenv("AWS_QUEUE_NAME")),
		DelaySeconds: aws.Int64(3),
	}
	send_resp, err := sqsSvc.SendMessage(send_params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("[Send message] \n%v \n\n", send_resp)

}

func Read() {

	chnMessages := make(chan *sqs.Message, 2)
	go PollMessages(chnMessages)

	for message := range chnMessages {
		handleMessage(message)
		SQSDeleteMessage(message, os.Getenv("AWS_QUEUE_NAME"))
	}

}

func PollMessages(chn chan<- *sqs.Message) {
	sqsSvc := SQSConnect()
	for {
		output, err := sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(os.Getenv("AWS_QUEUE_NAME")),
			MaxNumberOfMessages: aws.Int64(2),
			WaitTimeSeconds:     aws.Int64(15),
		})

		if err != nil {
			log.Println(err)
		}

		for _, message := range output.Messages {
			chn <- message
		}

	}

}

func SQSDeleteMessage(msg *sqs.Message, QueueUrl string) {
	sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(QueueUrl),
		ReceiptHandle: msg.ReceiptHandle,
	})
}

func handleMessage(msg *sqs.Message) {
	items := make(map[string]interface{})
	if err := json.Unmarshal([]byte(*msg.Body), &items); err != nil {
		log.Println(err)
	}
	log.Println(items["replyToken"])
	log.Println(items["messages"])
	if items["replyToken"] != nil {
		Reply(items["replyToken"].(string), items["messages"].(map[string]interface{}))

	}

}

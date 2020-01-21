package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

type Payload struct {
	Comment string
	Place   string
	Url     string
}

func HandleRequest(payload Payload) error {
	log.Printf("%v", payload)
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}

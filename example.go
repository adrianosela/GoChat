package main

import (
	"log"
	"net/http"

	"github.com/adrianosela/GoChat/chat"
)

func main() {
	svc := chat.NewChatService()
	err := http.ListenAndServe(":8080", svc)
	if err != nil {
		log.Fatal(err)
	}
}

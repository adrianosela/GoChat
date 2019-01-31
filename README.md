## GoChat! - A simple chat server for your Go APIs

This Project is inspired by the [chat example from the gorilla/websockets library.](https://github.com/gorilla/websocket/tree/master/examples/chat)

Currently the controller only broadcasts messages to all peers. Enabling messaging specific peers by ID is a work in progress

### Usage:

```
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
```


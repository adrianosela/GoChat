## GoChat! - A simple chat server for your Go APIs

This Project is inspired by the [chat example from the gorilla/websocket library.](https://github.com/gorilla/websocket/tree/master/examples/chat)

### WIP:
* Making a generic enough interface so that developers can adapt the chat service to their already existing APIs
* Adding end-to-end encryption in a way that the controller cannot see messages from peers (i.e. not SSL)
* Defining appropriate constants and defaults to the websocket connections and the controller's channels

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


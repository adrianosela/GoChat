package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// NewChatService returns the HTTP handler for the chat service
func NewChatService() http.Handler {
	c := NewController()
	go c.Start()

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/").HandlerFunc(serveHTML)
	r.Methods(http.MethodGet).Path("/ws").HandlerFunc(c.serveWS)
	return r
}

// serveHTML serves the "homepage"
func serveHTML(w http.ResponseWriter, r *http.Request) { w.Write([]byte(indexHTML)) }

// serveWS upgrades HTTP to websockets and creates a new peer for the request
func (c *Controller) serveWS(w http.ResponseWriter, r *http.Request) {
	// upgrade protocol to websockets connection
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// create new peer
	NewPeer(c, conn, make(chan []byte, 256), make(chan []byte, 256)).enroll()
}

/* We do this funny thing because this package is meant work as a Go library;
 * that is, the user should be able to use the plug-in chat server
 * after doing a go get and calling a function
 */
const indexHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
<title>Chat Service</title>
<script type="text/javascript">
window.onload = function () {
    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");
    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }
    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        conn.send(msg.value);
        msg.value = "";
        return false;
    };
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (event) {
            var messages = event.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                item.innerText = messages[i];
                appendLog(item);
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>this browser does not support websockets</b>";
        appendLog(item);
    }
};
</script>
<style type="text/css">
html {
    overflow: hidden;
}
body {
    overflow: hidden;
    padding: 0;
    margin: 0;
    width: 100%;
    height: 100%;
    background: gray;
}
#log {
    background: white;
    margin: 0;
    padding: 0.5em 0.5em 0.5em 0.5em;
    position: absolute;
    top: 0.5em;
    left: 0.5em;
    right: 0.5em;
    bottom: 3em;
    overflow: auto;
}
#form {
    padding: 0 0.5em 0 0.5em;
    margin: 0;
    position: absolute;
    bottom: 1em;
    left: 0px;
    width: 100%;
    overflow: hidden;
}
</style>
</head>
<body>
<div id="log"></div>
<form id="form">
    <input type="submit" value="Send" />
    <input type="text" id="msg" size="64"/>
</form>
</body>
</html>
`

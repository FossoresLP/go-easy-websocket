GoEasyWebsocket
===============
This package allows you to easily create web application servers with websocket support. Communication is processed using a text based structure consisting of a command and optionally data. The command/data has to be processed using handle functions for commands which may return a string submitted as a answer to the client. The data has to be processed in the handle functions.

A handler for `"open"` can be registered and will be called every time a connection is opened. Data will always be `nil`

Sample application using httprouter:
------------------------------------
```go
package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	ews "github.com/FossoresLP/go-easy-websocket"
)

func testWsHandler(data []byte) (resp []byte, err error) {
	resp = []byte("yourData: ") + data
	return
}

func main() {
	router := httprouter.New()
	ws := ews.NewHandler()
	ws.Handle("test", testWsHandler)
	router.GET("/ws", ws.UpgradeHandler)
}
```

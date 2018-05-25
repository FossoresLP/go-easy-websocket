GoEasyWebsocket
===============
This package allows you to easily create web application servers with websocket support. 
Communication is processed using a text based structure consisting of a command and optionally data. 

Client requests:
----------------
For client requests, handler functions can be registered for specific commands. They will recieve the raw data and may return a response.

A handler for `"open"` can be registered and will be called every time a connection is opened. The data submitted will always be the clients UUID.

Server push:
------------
The server can push commands and data to the clients using channels for which the clients have to register as listeners. The server may also push commands and data to specific clients whenever necessary using their UUID.

Please refer to GoDocs instead of the following example to better understand the package.
-----------------------------------------------------------------------------------------
> Sample application using httprouter:
> ------------------------------------
> ```go
> package main
> 
> import (
> 	"net/http"
> 	"github.com/julienschmidt/httprouter"
> 	ews "github.com/FossoresLP/go-easy-websocket"
> )
> 
> func testWsHandler(data []byte) (resp []byte, err error) {
> 	resp = []byte("yourData: ") + data
> 	return
> }
> 
> func main() {
> 	router := httprouter.New()
> 	ws := ews.NewHandler()
> 	ws.Handle("test", testWsHandler)
> 	router.GET("/ws", ws.UpgradeHandler)
> }
> ```

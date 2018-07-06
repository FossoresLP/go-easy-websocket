GoEasyWebsocket
===============

[![CircleCI](https://img.shields.io/circleci/project/github/FossoresLP/go-easy-websocket/master.svg?style=flat-square)](https://circleci.com/gh/FossoresLP/go-easy-websocket)

[![Coveralls](https://img.shields.io/coveralls/github/FossoresLP/go-easy-websocket/master.svg?style=flat-square)](https://coveralls.io/github/FossoresLP/go-easy-websocket)

[![Codacy](https://img.shields.io/codacy/grade/4d5e2ea1070046ab97d45df80ae09814.svg?style=flat-square)](https://www.codacy.com/app/FossoresLP/go-easy-websocket)


[![Licensed under: Boost Software License](https://img.shields.io/badge/style-BSL--1.0-red.svg?longCache=true&style=flat-square&label=License)](https://github.com/FossoresLP/go-easy-websocket/blob/master/LICENSE.md)
[![GoDoc](https://img.shields.io/badge/style-reference-blue.svg?longCache=true&style=flat-square&label=GoDoc)](https://godoc.org/github.com/FossoresLP/go-easy-websocket)

This package allows you to easily create web application servers with websocket support.
Communication is processed using a text based structure consisting of a command and optionally data.

Setup
-----

To create a new websocket handler call [NewHandler()](https://godoc.org/github.com/FossoresLP/go-easy-websocket#NewHandler).

Now you have to use [Handler.UpgradeHandler](https://godoc.org/github.com/FossoresLP/go-easy-websocket#Handler.#Handler.UpgradeHandler) to upgrade a normal HTTP request to websocket.

```go
handler := websocket.NewHandler()

http.HandleFunc("/ws", handler.UpgradeHandler)

log.Fatal(http.ListenAndServe(":8080", nil))
```

Client requests
---------------

You can handle client requests by registering [handle functions](https://godoc.org/github.com/FossoresLP/go-easy-websocket#HandleFunc) for specific command strings using [Handler.Handle](https://godoc.org/github.com/FossoresLP/go-easy-websocket#Handler.Handle)

A handler for `open` can be registered and will be called every time a connection is opened. The message will be the clients session ID.

```go
handler.Handle("open", func(msg []byte, authToken string) *Message {
	// Parse message here
	sessionid = uuid.Parse(msg)
	// Validate auth token if necessary
	return NewMessage("welcome", []byte("Hello " string(msg)))
	// This will write "Hello xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" to the client using the command "welcome"
})
```

Server push
-----------

The server can push commands and data to the clients using channels for which the clients have to register as listeners. The server may also push commands and data to specific clients whenever necessary using their session ID.

```go
handler.RegisterListenChannel("test")

handler.WriteToChannel("test", NewMessage("thanks", []byte("Thank you for listening!")))

handler.WriteToClient(sessionid, NewMessage("direct", []byte("This message is only sent to a single client")))
```

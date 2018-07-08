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

The server can push commands and data to the clients using channels for which the clients have to register as listeners.

Channels can have validation functions attached on setup to restrict which clients will be allowed to register as listeners.
These validation functions will be called whenever a client tries to register as a listener with the auth token that client submitted when connecting.
A return value of `nil` will be considered a successful validation while any error will be considered a validation failure and therefore prevent the client from registering as a listener. The errors will not be relayed to the client to improve security. Instead a generic error message will be sent.

The server may also push commands and data to specific clients whenever necessary using their session ID.

```go
handler.RegisterListenChannel("test", nil) // Anyone can register as a listener

handler.RegisterListenChannel("restricted", func(t string) error {
	if t != "valid" {
		return errors.New("Invalid token")
	}
	return nil
})

handler.WriteToChannel("test", NewMessage("thanks", []byte("Thank you for listening!")))

handler.WriteToClient(sessionid, NewMessage("direct", []byte("This message is only sent to a single client")))
```

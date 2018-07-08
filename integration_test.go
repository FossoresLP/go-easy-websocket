package websocket_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	ws "github.com/fossoreslp/go-easy-websocket"
	wsc "github.com/gorilla/websocket"
)

func serverRoutine(t *testing.T) {
	h := ws.NewHandler()
	h.ValidateFunction = func(s string) error {
		if s != "valid" {
			return errors.New("Failed to authenticate")
		}
		return nil
	}
	h.Handle("open", func(in []byte, _ string) *ws.Message {
		msg, _ := ws.NewMessage("hello", []byte("Connected"))
		return msg
	})
	h.Handle("testChannel", func(_ []byte, _ string) *ws.Message {
		msg, _ := ws.NewMessage("channel", []byte("test"))
		h.WriteToChannel("test", msg)
		return nil
	})
	h.Handle("testResponse", func(_ []byte, _ string) *ws.Message {
		msg, _ := ws.NewMessage("response", []byte("sent"))
		return msg
	})
	h.RegisterListenChannel("test", nil)
	h.RegisterListenChannel("validate", func(t string) error {
		if t != "channel_valid" {
			return errors.New("Not permitted")
		}
		return nil
	})
	http.HandleFunc("/", h.UpgradeHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		t.Fatalf("HTTP server failed: %s", err.Error())
	}
}

func initClient() (*wsc.Conn, error) {
	dialer := wsc.DefaultDialer
	dialer.Subprotocols = []string{"cmd.fossores.de"}
	client, resp, err := dialer.Dial("ws://localhost:8080/", http.Header{"Sec-WebSocket-Protocol": []string{"cmd.fossores.de"}, "Cookie": []string{"auth=valid"}})
	if err != nil {
		b := make([]byte, 64)
		resp.Body.Read(b)
		return client, fmt.Errorf("Failed to open websocket connection: %s with response: %+v, %q, %+v", err.Error(), *resp, b, resp.Request)
	}
	if resp.Header.Get("Sec-WebSocket-Protocol") != "cmd.fossores.de" {
		return client, fmt.Errorf("Invalid websocket protocol: %q", resp.Header.Get("Sec-WebSocket-Protocol"))
	}
	return client, nil
}

func Test_Websocket(t *testing.T) {
	go serverRoutine(t) // nolint: errcheck

	client, err := initClient()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer client.Close() // nolint: errcheck
	// Receive open message
	_, msg, err := client.ReadMessage()
	if err != nil {
		t.Errorf("Failed to receive message: %s", err.Error())
	}
	t.Log(msg)
	// Register as listener
	err = client.WriteMessage(wsc.TextMessage, []byte("listen: test"))
	if err != nil {
		t.Errorf("Failed to send message: %s", err.Error())
	}
	// Ask for channel test
	err = client.WriteMessage(wsc.TextMessage, []byte("testChannel: "))
	if err != nil {
		t.Errorf("Failed to send message: %s", err.Error())
	}
	// Receive channel test
	_, msg, err = client.ReadMessage()
	if err != nil {
		t.Errorf("Failed to receive message: %s", err.Error())
	}
	t.Log(msg)
	// Send invalid command
	err = client.WriteMessage(wsc.TextMessage, []byte("invalid: command"))
	if err != nil {
		t.Errorf("Failed to send message: %s", err.Error())
	}
	// Receive response
	_, msg, err = client.ReadMessage()
	if err != nil {
		t.Errorf("Failed to receive message: %s", err.Error())
	}
	t.Log(msg)
	// Request response test
	err = client.WriteMessage(wsc.TextMessage, []byte("testResponse: "))
	if err != nil {
		t.Errorf("Failed to send message: %s", err.Error())
	}
	// Receive response test
	_, msg, err = client.ReadMessage()
	if err != nil {
		t.Errorf("Failed to receive message: %s", err.Error())
	}
	t.Log(msg)
	// Register for invalid channel
	err = client.WriteMessage(wsc.TextMessage, []byte("listen: invalid"))
	if err != nil {
		t.Errorf("Failed to send message: %s", err.Error())
	}
	// Receive response
	_, msg, err = client.ReadMessage()
	if err != nil {
		t.Errorf("Failed to receive message: %s", err.Error())
	}
	t.Log(msg)
	// Register for channel without authentication
	err = client.WriteMessage(wsc.TextMessage, []byte("listen: validate"))
	if err != nil {
		t.Errorf("Failed to send message: %s", err.Error())
	}
	// Receive response
	_, msg, err = client.ReadMessage()
	if err != nil {
		t.Errorf("Failed to receive message: %s", err.Error())
	}
	t.Log(msg)
}

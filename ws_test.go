package websocket

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type TestHTTPResponse struct {
	header  http.Header
	content *TestHTTPResponseContent
}

type TestHTTPResponseContent struct {
	status  int
	written []byte
}

func (resp TestHTTPResponse) Write(c []byte) (int, error) {
	resp.content.written = c
	return len(c), nil
}
func (resp TestHTTPResponse) WriteHeader(c int) {
	resp.content.status = c
}
func (resp TestHTTPResponse) Header() http.Header {
	return resp.header
}
func (resp TestHTTPResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return TestNetConn(true), bufio.NewReadWriter(bufio.NewReader(bytes.NewBuffer(make([]byte, 10))), bufio.NewWriter(bytes.NewBuffer(make([]byte, 10)))), nil
}

type TestNetConn bool

type TestNetAddr bool

func (c TestNetAddr) Network() string {
	return ""
}
func (c TestNetAddr) String() string {
	return ""
}

func (c TestNetConn) Read(b []byte) (int, error) {
	return len(b), nil
}
func (c TestNetConn) Write(b []byte) (int, error) {
	return len(b), nil
}
func (c TestNetConn) Close() error {
	return nil
}
func (c TestNetConn) LocalAddr() net.Addr {
	return TestNetAddr(true)
}
func (c TestNetConn) RemoteAddr() net.Addr {
	return TestNetAddr(true)
}
func (c TestNetConn) SetDeadline(t time.Time) error {
	return nil
}
func (c TestNetConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (c TestNetConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestHandler_UpgradeHandler(t *testing.T) {
	t.Run("Testing the test", func(t *testing.T) {
		r := TestHTTPResponse{content: &TestHTTPResponseContent{0, nil}}
		r.WriteHeader(200)
		l, _ := r.Write([]byte("Test"))
		if l != 4 {
			t.Error("len() failed")
		}
		if r.content.status == 0 {
			t.Error("Status has not been set")
		}
		if r.content.written == nil {
			t.Error("Written has not been set")
		}
	})
	t.Run("Normal", func(t *testing.T) {
		h := NewHandler()
		h.ValidateFunction = func(c string) error {
			if c == "valid" {
				return nil
			}
			return errors.New("invalid")
		}
		req, err := http.NewRequest("GET", "wss://example.com/ws", http.NoBody)
		if err != nil {
			t.Errorf("Could not create new request, tests cannot be run: %s", err.Error())
		}
		req.AddCookie(&http.Cookie{Name: "auth", Value: "valid"})
		req.Header.Add("Connection", "upgrade")
		req.Header.Add("Upgrade", "websocket")
		req.Header.Add("Sec-Websocket-Version", "13")
		req.Header.Add("Sec-Websocket-Key", "test")
		resp := TestHTTPResponse{content: &TestHTTPResponseContent{0, nil}}
		h.UpgradeHandler(resp, req)
		if resp.content.status != 0 {
			t.Errorf("HTTP status should not be set on normal connection but is %d", resp.content.status)
		}
		if resp.content.written != nil {
			t.Errorf("HTTP body should not be set on normal connection but is %q", resp.content.written)
		}
	})
	t.Run("NoAuthCookie", func(t *testing.T) {
		h := NewHandler()
		h.ValidateFunction = func(c string) error {
			if c == "valid" {
				return nil
			}
			return errors.New("invalid")
		}
		req, err := http.NewRequest("GET", "wss://example.com/ws", http.NoBody)
		if err != nil {
			t.Errorf("Could not create new request, tests cannot be run: %s", err.Error())
		}
		resp := TestHTTPResponse{content: &TestHTTPResponseContent{0, nil}}
		h.UpgradeHandler(resp, req)
		if resp.content.status != 403 {
			t.Errorf("HTTP status should be 403 but is %d", resp.content.status)
		}
		if !reflect.DeepEqual(resp.content.written, []byte("Authentication failed\n")) {
			t.Errorf("HTTP body should be \"Authentication failed\" but is %q", resp.content.written)
		}
	})
	t.Run("InvalidAuthCookie", func(t *testing.T) {
		h := NewHandler()
		h.ValidateFunction = func(c string) error {
			if c == "valid" {
				return nil
			}
			return errors.New("invalid")
		}
		req, err := http.NewRequest("GET", "wss://example.com/ws", http.NoBody)
		if err != nil {
			t.Errorf("Could not create new request, tests cannot be run: %s", err.Error())
		}
		req.AddCookie(&http.Cookie{Name: "auth", Value: "invalid"})
		resp := TestHTTPResponse{content: &TestHTTPResponseContent{0, nil}}
		h.UpgradeHandler(resp, req)
		if resp.content.status != 403 {
			t.Errorf("HTTP status should be 403 but is %d", resp.content.status)
		}
		if !reflect.DeepEqual(resp.content.written, []byte("Authentication failed\n")) {
			t.Errorf("HTTP body should be \"Authentication failed\" but is %q", resp.content.written)
		}
	})
	t.Run("NoConnectionUpgrade", func(t *testing.T) {
		h := NewHandler()
		h.ValidateFunction = func(c string) error {
			if c == "valid" {
				return nil
			}
			return errors.New("invalid")
		}
		req, err := http.NewRequest("GET", "wss://example.com/ws", http.NoBody)
		if err != nil {
			t.Errorf("Could not create new request, tests cannot be run: %s", err.Error())
		}
		req.AddCookie(&http.Cookie{Name: "auth", Value: "valid"})
		resp := TestHTTPResponse{content: &TestHTTPResponseContent{0, nil}, header: http.Header{}}
		h.UpgradeHandler(resp, req)
		if resp.content.status != 426 {
			t.Errorf("HTTP status should be 426 but is %d", resp.content.status)
		}
		if resp.header.Get("Upgrade") != "WebSocket" {
			t.Errorf("HTTP header `Upgrade` should be `WebSocket` but is %q", resp.header.Get("Upgrade"))
		}
	})
	t.Run("NoUpgradeWebSocket", func(t *testing.T) {
		h := NewHandler()
		h.ValidateFunction = func(c string) error {
			if c == "valid" {
				return nil
			}
			return errors.New("invalid")
		}
		req, err := http.NewRequest("GET", "wss://example.com/ws", http.NoBody)
		if err != nil {
			t.Errorf("Could not create new request, tests cannot be run: %s", err.Error())
		}
		req.AddCookie(&http.Cookie{Name: "auth", Value: "valid"})
		req.Header.Add("Connection", "upgrade")
		resp := TestHTTPResponse{content: &TestHTTPResponseContent{0, nil}, header: http.Header{}}
		h.UpgradeHandler(resp, req)
		if resp.content.status != 426 {
			t.Errorf("HTTP status should be 426 but is %d", resp.content.status)
		}
		if resp.header.Get("Upgrade") != "WebSocket" {
			t.Errorf("HTTP header `Upgrade` should be `WebSocket` but is %q", resp.header.Get("Upgrade"))
		}
	})
	t.Run("NoSecWebSocketVersion13", func(t *testing.T) {
		h := NewHandler()
		h.ValidateFunction = func(c string) error {
			if c == "valid" {
				return nil
			}
			return errors.New("invalid")
		}
		req, err := http.NewRequest("GET", "wss://example.com/ws", http.NoBody)
		if err != nil {
			t.Errorf("Could not create new request, tests cannot be run: %s", err.Error())
		}
		req.AddCookie(&http.Cookie{Name: "auth", Value: "valid"})
		req.Header.Add("Connection", "upgrade")
		req.Header.Add("Upgrade", "websocket")
		resp := TestHTTPResponse{content: &TestHTTPResponseContent{0, nil}, header: http.Header{}}
		h.UpgradeHandler(resp, req)
		if resp.content.status != 426 {
			t.Errorf("HTTP status should be 426 but is %d", resp.content.status)
		}
		if resp.header.Get("Upgrade") != "WebSocket" {
			t.Errorf("HTTP header `Upgrade` should be `WebSocket` but is %q", resp.header.Get("Upgrade"))
		}
	})
	t.Run("NoSecWebSocketKey", func(t *testing.T) {
		h := NewHandler()
		h.ValidateFunction = func(c string) error {
			if c == "valid" {
				return nil
			}
			return errors.New("invalid")
		}
		req, err := http.NewRequest("GET", "wss://example.com/ws", http.NoBody)
		if err != nil {
			t.Errorf("Could not create new request, tests cannot be run: %s", err.Error())
		}
		req.AddCookie(&http.Cookie{Name: "auth", Value: "valid"})
		req.Header.Add("Connection", "upgrade")
		req.Header.Add("Upgrade", "websocket")
		req.Header.Add("Sec-Websocket-Version", "13")
		resp := TestHTTPResponse{content: &TestHTTPResponseContent{0, nil}, header: http.Header{}}
		h.UpgradeHandler(resp, req)
		if resp.content.status != 426 {
			t.Errorf("HTTP status should be 426 but is %d", resp.content.status)
		}
		if resp.header.Get("Upgrade") != "WebSocket" {
			t.Errorf("HTTP header `Upgrade` should be `WebSocket` but is %q", resp.header.Get("Upgrade"))
		}
	})
}

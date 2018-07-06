package websocket

import (
	"reflect"
	"testing"
)

func TestHandler_Handle(t *testing.T) {
	handler := NewHandler()
	handler.Handle("duplicate", func(b []byte, s string) *Message {
		return &Message{}
	})
	type args struct {
		cmd    string
		action HandleFunc
	}
	tests := []struct {
		name    string
		h       *Handler
		args    args
		wantErr bool
	}{
		{"Normal", NewHandler(), args{"test", func(b []byte, s string) *Message {
			return &Message{}
		}}, false},
		{"CommandOver255Characters", NewHandler(), args{"This command goes on for more than 255 characters which is not supported to keep the message size down. The limit of 255 characters has been chosen because we add a colon after the command and therefore effectively use 256 characters for the command. This limit should never be a problem unless you try to use the command to transmit data which is not recommended.", func(b []byte, s string) *Message {
			return &Message{}
		}}, true},
		{"CommandContainsColon", NewHandler(), args{"test: with colon", func(b []byte, s string) *Message {
			return &Message{}
		}}, true},
		{"CommandWebSocketIsReserved", NewHandler(), args{"websocket", func(b []byte, s string) *Message {
			return &Message{}
		}}, true},
		{"Command", handler, args{"duplicate", func(b []byte, s string) *Message {
			return &Message{}
		}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.Handle(tt.args.cmd, tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("Handler.Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseMessage(t *testing.T) {
	type args struct {
		msg []byte
	}
	tests := []struct {
		name    string
		args    args
		wantMsg *Message
	}{
		{"Normal", args{[]byte("command: message")}, &Message{[]byte("command"), []byte("message")}},
		{"NoColon", args{[]byte("command message")}, &Message{}},
		{"EmptyCommand", args{[]byte(": message")}, &Message{}},
		{"CommandExceedsLengthLimit", args{[]byte("This command goes on for more than 255 characters which is not supported to keep the message size down. The limit of 255 characters has been chosen because we add a colon after the command and therefore effectively use 256 characters for the command. This limit should never be a problem unless you try to use the command to transmit data which is not recommended.: message")}, &Message{}},
		{"CommandOnly", args{[]byte("command: ")}, &Message{[]byte("command"), nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsg := parseMessage(tt.args.msg)
			if !reflect.DeepEqual(gotMsg, tt.wantMsg) {
				t.Errorf("parseMessage() gotMsg = %v, want %v", gotMsg, tt.wantMsg)
			}
		})
	}
}

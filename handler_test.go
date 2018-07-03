package websocket

import (
	"reflect"
	"testing"
)

func TestHandler_Handle(t *testing.T) {
	handler := NewHandler()
	handler.Handle("duplicate", func(b []byte, s string) (*Message, error) {
		return &Message{}, nil
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
		{"Normal", NewHandler(), args{"test", func(b []byte, s string) (*Message, error) {
			return &Message{}, nil
		}}, false},
		{"CommandOver255Characters", NewHandler(), args{"This command goes on for more than 255 characters which is not supported to keep the message size down. The limit of 255 characters has been chosen because we add a colon after the command and therefore effectively use 256 characters for the command. This limit should never be a problem unless you try to use the command to transmit data which is not recommended.", func(b []byte, s string) (*Message, error) {
			return &Message{}, nil
		}}, true},
		{"CommandContainsColon", NewHandler(), args{"test: with colon", func(b []byte, s string) (*Message, error) {
			return &Message{}, nil
		}}, true},
		{"CommandWebSocketIsReserved", NewHandler(), args{"websocket", func(b []byte, s string) (*Message, error) {
			return &Message{}, nil
		}}, true},
		{"Command", handler, args{"duplicate", func(b []byte, s string) (*Message, error) {
			return &Message{}, nil
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

func Test_parseMsgByte(t *testing.T) {
	type args struct {
		msg []byte
	}
	tests := []struct {
		name     string
		args     args
		wantCmd  string
		wantData []byte
		wantErr  bool
	}{
		{"Normal", args{[]byte("command: message")}, "command", []byte("message"), false},
		{"NoColon", args{[]byte("command message")}, "", nil, true},
		{"EmptyCommand", args{[]byte(": message")}, "", []byte("message"), true},
		{"CommandExceedsLengthLimit", args{[]byte("This command goes on for more than 255 characters which is not supported to keep the message size down. The limit of 255 characters has been chosen because we add a colon after the command and therefore effectively use 256 characters for the command. This limit should never be a problem unless you try to use the command to transmit data which is not recommended.: message")}, "", nil, true},
		{"CommandOnly", args{[]byte("command: ")}, "command", []byte{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotData, err := parseMsgByte(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMsgByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCmd != tt.wantCmd {
				t.Errorf("parseMsgByte() gotCmd = %v, want %v", gotCmd, tt.wantCmd)
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("parseMsgByte() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

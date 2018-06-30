package websocket

import (
	"reflect"
	"testing"
)

func TestNewMessage(t *testing.T) {
	type args struct {
		cmd  string
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Message
		wantErr bool
	}{
		{"Normal", args{"test", nil}, Message{[]byte("test"), nil}, false},
		{"CommandTooLong", args{"This command goes on for more than 255 characters which is not supported to keep the message size down. The limit of 255 characters has been chosen because we add a colon after the command and therefore effectively use 256 characters for the command. This limit should never be a problem unless you try to use the command to transmit data which is not recommended.", nil}, Message{}, true},
		{"CommandEmpty", args{"", nil}, Message{}, true},
		{"CommandWithColon", args{"testing: colons", nil}, Message{}, true},
		{"CommandWebSocket", args{"websocket", nil}, Message{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMessage(tt.args.cmd, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("NewMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

package websocket

import (
	"reflect"
	"testing"

	"github.com/fossoreslp/go-uuid-v4"
)

func TestHandler_writeToClient(t *testing.T) {
	h := NewHandler()
	sessionID, err := uuid.New()
	if err != nil {
		t.Fatalf("Failed to generate session ID for testing: %s", err.Error())
	}
	h.writeChannels[sessionID] = make(chan []byte, 8)
	t.Run("Normal", func(t *testing.T) {
		err = h.writeToClient(sessionID, []byte("cmd"), []byte("content"))
		if err != nil {
			t.Errorf("Write failed unexpectedly: %s", err.Error())
		}
		msg := <-h.writeChannels[sessionID]
		if !reflect.DeepEqual(msg, []byte("cmd: content")) {
			t.Errorf("Invalid message. Should be %q but is %q.", []byte("cmd: content"), msg)
		}
	})
	t.Run("InvalidSessionID", func(t *testing.T) {
		randomID, err := uuid.New()
		if err != nil {
			t.Fatalf("Failed to generate random ID for testing: %s", err.Error())
		}
		err = h.writeToClient(randomID, []byte("cmd"), []byte("content"))
		if err == nil {
			t.Error("Write to invalid ID should fail")
		}
	})
}

func TestHandler_WriteToClient(t *testing.T) {
	h := NewHandler()
	sessionID, err := uuid.New()
	if err != nil {
		t.Fatalf("Failed to generate session ID for testing: %s", err.Error())
	}
	h.writeChannels[sessionID] = make(chan []byte, 8)
	t.Run("Normal", func(t *testing.T) {
		err = h.WriteToClient(sessionID, &Message{[]byte("cmd"), []byte("content")})
		if err != nil {
			t.Errorf("Write failed unexpectedly: %s", err.Error())
		}
		msg := <-h.writeChannels[sessionID]
		if !reflect.DeepEqual(msg, []byte("cmd: content")) {
			t.Errorf("Invalid message. Should be %q but is %q.", []byte("cmd: content"), msg)
		}
	})
	t.Run("CommandEmpty", func(t *testing.T) {
		err = h.WriteToClient(sessionID, &Message{nil, []byte("content")})
		if err == nil || err.Error() != "command may not be empty" {
			t.Error("Empty command not detected")
		}
	})
	t.Run("CommandTooLong", func(t *testing.T) {
		err = h.WriteToClient(sessionID, &Message{[]byte("This command goes on for more than 255 characters which is not supported to keep the message size down. The limit of 255 characters has been chosen because we add a colon after the command and therefore effectively use 256 characters for the command. This limit should never be a problem unless you try to use the command to transmit data which is not recommended."), []byte("content")})
		if err == nil || err.Error() != "command may not be longer than 255 characters" {
			t.Error("Command with more than 255 characters not detected")
		}
	})
}

func TestHandler_channelRoutine(t *testing.T) {
	h := NewHandler()
	sessionID, err := uuid.New()
	if err != nil {
		t.Fatalf("Failed to generate session ID for testing: %s", err.Error())
	}
	randomID, err := uuid.New()
	if err != nil {
		t.Fatalf("Failed to generate random ID for testing: %s", err.Error())
	}
	h.writeChannels[sessionID] = make(chan []byte, 8)
	h.channels["test"] = &channel{make(chan *Message, 2), []uuid.UUID{sessionID, randomID}}
	go h.channelRoutine("test")
	h.channels["test"].send <- &Message{[]byte("cmd"), []byte("content")}
	msg := <-h.writeChannels[sessionID]
	if !reflect.DeepEqual(msg, []byte("cmd: content")) {
		t.Errorf("Invalid message. Should be %q but is %q.", []byte("cmd: content"), msg)
	}
}

func TestHandler_WriteToChannel(t *testing.T) {
	h := NewHandler()
	h.channels["test"] = &channel{make(chan *Message, 8), []uuid.UUID{}}
	type args struct {
		channel string
		msg     *Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Normal", args{"test", &Message{[]byte("cmd"), []byte("content")}}, false},
		{"CommandEmpty", args{"test", &Message{nil, []byte("content")}}, true},
		{"CommandTooLong", args{"test", &Message{[]byte("This command goes on for more than 255 characters which is not supported to keep the message size down. The limit of 255 characters has been chosen because we add a colon after the command and therefore effectively use 256 characters for the command. This limit should never be a problem unless you try to use the command to transmit data which is not recommended."), []byte("content")}}, true},
		{"ChannelNotFound", args{"imaginary", &Message{[]byte("cmd"), []byte("content")}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := h.WriteToChannel(tt.args.channel, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Handler.WriteToChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				msg := <-h.channels["test"].send
				if !reflect.DeepEqual(msg, &Message{[]byte("cmd"), []byte("content")}) {
					t.Errorf("Message does not match expected (cmd: content): %+v", msg)
				}
			}
		})
	}
}

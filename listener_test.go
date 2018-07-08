package websocket

import (
	"testing"

	"github.com/fossoreslp/go-uuid-v4"
)

func TestHandler_RegisterListenChannel(t *testing.T) {
	handler := NewHandler()
	type args struct {
		name         string
		validateFunc func(string) error
	}
	tests := []struct {
		name    string
		h       *Handler
		args    args
		wantErr bool
	}{
		{"Normal", handler, args{"test", nil}, false},
		{"AlreadyExists", handler, args{"test", nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.RegisterListenChannel(tt.args.name, tt.args.validateFunc); (err != nil) != tt.wantErr {
				t.Errorf("Handler.RegisterListenChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_registerAsListener(t *testing.T) {
	handler := NewHandler()
	id := uuid.UUID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	if err := handler.RegisterListenChannel("test", nil); err != nil {
		t.Fatalf("Could not register listen channel: %s", err.Error())
	}
	type args struct {
		id   uuid.UUID
		name string
	}
	tests := []struct {
		name    string
		h       *Handler
		args    args
		wantErr bool
	}{
		{"Normal", handler, args{id, "test"}, false},
		{"InvalidChannel", handler, args{id, "default"}, true},
		{"AlreadyListening", handler, args{id, "test"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.registerAsListener(tt.args.id, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Handler.registerAsListener() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if len(handler.channels["test"].listeners) == 0 || handler.channels["test"].listeners[0] != id {
		t.Error("Handler.registerAsListener() failed to register id")
	}

}

func TestHandler_unregisterAsListener(t *testing.T) {
	// Initialize values
	handler := NewHandler()
	id := uuid.UUID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	altID := uuid.UUID{0xF, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	if err := handler.RegisterListenChannel("test", nil); err != nil {
		t.Fatalf("Could not register listen channel: %s", err.Error())
	}
	// Add ID as only ID in list
	handler.channels["test"].listeners = append(handler.channels["test"].listeners, id)
	// Try to unregister
	if err := handler.unregisterAsListener(id, "test"); err != nil {
		t.Errorf("Handler.unregisterAsListener() failed: %s", err.Error())
	}
	// Check if ID is properly removed
	if len(handler.channels["test"].listeners) > 0 {
		t.Error("Handler.unregisterAsListener() failed to unregister id")
	}
	// Try to unregister from invalid channel and check for proper error
	if err := handler.unregisterAsListener(id, "default"); err == nil {
		t.Error("Handler.unregisterAsListener() failed to detect invalid channel")
	}
	// Add ID and altID to list
	handler.channels["test"].listeners = append(handler.channels["test"].listeners, id, altID)
	// Try to unregister
	if err := handler.unregisterAsListener(id, "test"); err != nil {
		t.Errorf("Handler.unregisterAsListener() failed: %s", err.Error())
	}
	// Check if ID is properly removed
	if len(handler.channels["test"].listeners) != 1 || handler.channels["test"].listeners[0] == id {
		t.Error("Handler.unregisterAsListener() failed to unregister id")
	}
	// Add ID to end of list
	handler.channels["test"].listeners = append(handler.channels["test"].listeners, id)
	// Try to unregister
	if err := handler.unregisterAsListener(id, "test"); err != nil {
		t.Errorf("Handler.unregisterAsListener() failed: %s", err.Error())
	}
	// Check if ID is properly removed
	if len(handler.channels["test"].listeners) != 1 || handler.channels["test"].listeners[0] == id {
		t.Error("Handler.unregisterAsListener() failed to unregister id")
	}
	// Add ID and altID to list
	handler.channels["test"].listeners = append(handler.channels["test"].listeners, id, altID)
	// Try to unregister
	if err := handler.unregisterAsListener(id, "test"); err != nil {
		t.Errorf("Handler.unregisterAsListener() failed: %s", err.Error())
	}
	// Check if ID is properly removed
	if len(handler.channels["test"].listeners) != 2 || handler.channels["test"].listeners[0] == id || handler.channels["test"].listeners[1] == id {
		t.Error("Handler.unregisterAsListener() failed to unregister id")
	}
}

func TestHandler_unregisterListener(t *testing.T) {
	handler := NewHandler()
	id := uuid.UUID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	if err := handler.RegisterListenChannel("test1", nil); err != nil {
		t.Fatalf("Could not register listen channel: %s", err.Error())
	}
	handler.channels["test1"].listeners = append(handler.channels["test1"].listeners, id)
	if err := handler.RegisterListenChannel("test2", nil); err != nil {
		t.Fatalf("Could not register listen channel: %s", err.Error())
	}
	if err := handler.RegisterListenChannel("test3", nil); err != nil {
		t.Fatalf("Could not register listen channel: %s", err.Error())
	}
	handler.channels["test3"].listeners = append(handler.channels["test3"].listeners, id)
	if err := handler.RegisterListenChannel("test4", nil); err != nil {
		t.Fatalf("Could not register listen channel: %s", err.Error())
	}
	handler.unregisterListener(id)
	if len(handler.channels["test1"].listeners) > 0 || len(handler.channels["test3"].listeners) > 0 {
		t.Error("Handler.unregisterListener() failed to unregister id from all channels")
	}
}

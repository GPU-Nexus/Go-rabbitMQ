package main

import (
	"bytes"
	"go_rabbitMQ/messages"
	"os"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestMain(t *testing.T) {
	// Set up environment variables for testing
	err := os.Setenv("RABBITMQ_PASSWORD", "guest")
	if err != nil {
		t.Fatalf("Error setting environment variable: %s", err)
	}

	// Create a message
	msg := &messages.MyMessage{
		Id:      "test-id",
		Content: "test content",
	}

	// Marshal the message to Protobuf
	body, err := proto.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %s", err)
	}

	// Check if the message is correctly marshaled
	expected := []byte{10, 7 ,116, 101, 115, 116, 45, 105, 100, 18, 12, 116, 101, 115, 116, 32, 99, 111, 110, 116, 101, 110, 116} // Replace with the expected byte slice of the marshaled message
	if !bytes.Equal(body, expected) {
		t.Errorf("Expected %v, but got %v", expected, body)
	}

	// Clean up environment variables after test
	err = os.Unsetenv("RABBITMQ_PASSWORD")
	if err != nil {
		t.Fatalf("Error unsetting environment variable: %s", err)
	}
}

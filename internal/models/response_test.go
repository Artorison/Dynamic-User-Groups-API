package models

import (
	"fmt"
	"testing"
)

func TestOK(t *testing.T) {
	message := "test message"
	resp := OK(message)

	if resp.Message != message {
		t.Errorf("Expected 'test message', got '%s'", resp.Message)
	}

	if resp.Error != "" {
		t.Errorf("Expected empty error, got '%s'", resp.Error)
	}
}

func TestResponseErr(t *testing.T) {
	result := ResponseErr("error occurred")
	if result.Message != "error occurred" {
		t.Errorf("Expected message 'error occurred', got '%s'", result.Message)
	}
	if result.Error != "" {
		t.Errorf("Expected empty error, got '%s'", result.Error)
	}

	testErr := fmt.Errorf("test error")
	result = ResponseErr("error occurred", testErr)
	if result.Message != "error occurred" {
		t.Errorf("Expected message 'error occurred', got '%s'", result.Message)
	}
	if result.Error != testErr.Error() {
		t.Errorf("Expected error '%s', got '%s'", testErr.Error(), result.Error)
	}
}

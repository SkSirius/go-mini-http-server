package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHelloHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rr := httptest.NewRecorder()

	HelloHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	expected := `{"message":"Hello, world!"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("Unexpected body: got %s, want %s", rr.Body.String(), expected)
	}
}

func TestTimeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	rr := httptest.NewRecorder()

	TimeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("Invalid JSON response: %v", err)
	}

	if _, err := time.Parse(time.RFC3339, resp["time"]); err != nil {
		t.Errorf("Invalid time format: %v", err)
	}
}

func TestEchoHandler_ValidJSON(t *testing.T) {
	payload := `{"message":"hello"}`
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	EchoHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("Invalid JSON response: %v", err)
	}
	if resp["echo"] != "hello" {
		t.Errorf("Unexpected echo response: %s", resp["echo"])
	}
}

func TestEchoHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/echo", nil)
	rr := httptest.NewRecorder()

	EchoHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

func TestEchoHandler_BadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	EchoHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

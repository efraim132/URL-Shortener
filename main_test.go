package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRedirect(t *testing.T) {
	// Setup
	theWholeStore.data["/test"] = "https://example.com"

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handleRedirect(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected status %d, got %d", http.StatusFound, w.Code)
	}

	location := w.Header().Get("Location")
	if location != "https://example.com" {
		t.Errorf("Expected location %s, got %s", "https://example.com", location)
	}
}

func TestHandleRedirectNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	handleRedirect(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandlePostURL(t *testing.T) {
	entry := entry{Short: "gh", Long: "https://github.com"}
	body, _ := json.Marshal(entry)

	req := httptest.NewRequest("POST", "/urls", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handlePostURL(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Verify the URL was stored
	theWholeStore.mutex.RLock()
	stored := theWholeStore.data["/gh"]
	theWholeStore.mutex.RUnlock()

	if stored != "https://github.com" {
		t.Errorf("Expected stored URL %s, got %s", "https://github.com", stored)
	}
}

func TestHandlePostURLInvalidMethod(t *testing.T) {
	req := httptest.NewRequest("GET", "/urls", nil)
	w := httptest.NewRecorder()

	handlePostURL(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

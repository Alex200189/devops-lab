package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != `{"status":"ok"}` {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestCreateTasksEmptyTitle(t *testing.T) {
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":""}`))
	w := httptest.NewRecorder()

	createTasksHandler(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateTasksInvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`not json`))
	w := httptest.NewRecorder()

	createTasksHandler(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

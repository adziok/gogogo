package common

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRenderErrDoesNotExposeInternalError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/feature-flag", nil)
	recorder := httptest.NewRecorder()

	RenderErr(recorder, req, http.StatusInternalServerError, errors.New("postgres password leaked"), "Internal server error")

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", contentType)
	}

	var response ErrResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != http.StatusText(http.StatusInternalServerError) {
		t.Errorf("status = %q, want %q", response.Status, http.StatusText(http.StatusInternalServerError))
	}
	if response.Error != "Internal server error" {
		t.Errorf("error = %q, want %q", response.Error, "Internal server error")
	}
	if strings.Contains(recorder.Body.String(), "postgres password leaked") {
		t.Fatal("internal error was exposed in the response")
	}
}

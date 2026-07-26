package common

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func RenderErr(w http.ResponseWriter, r *http.Request, statusCode int, err error, message string) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if statusCode >= http.StatusInternalServerError {
			slog.Error(message, "error", err, "method", r.Method, "path", r.URL.Path)
		} else {
			slog.Warn(message, "error", err, "method", r.Method, "path", r.URL.Path)
		}
	}

	errBody := ErrResponse{
		Status: http.StatusText(statusCode),
		Error:  message,
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(errBody); err != nil {
		slog.Error("Failed to encode error response", "error", err, "method", r.Method, "path", r.URL.Path)
	}
}

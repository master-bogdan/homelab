// Package utils contains common utils
package utils

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
}

func WriteResponseError(w http.ResponseWriter, status int, msg string) {
	WriteResponse(w, status, map[string]string{"error": msg})
}

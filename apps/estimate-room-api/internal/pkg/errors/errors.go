package errors

import (
	"encoding/json"
	"net/http"
)

type ErrorItem struct {
	Detail  string `json:"detail"`
	Pointer string `json:"pointer"`
}

type Problem struct {
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Status   int         `json:"status"`
	Detail   string      `json:"detail"`
	Instance string      `json:"instance"`
	Errors   []ErrorItem `json:"errors"`
}

func Write(w http.ResponseWriter, problem Problem) {
	if problem.Errors == nil {
		problem.Errors = []ErrorItem{}
	}
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.WriteHeader(problem.Status)
	_ = json.NewEncoder(w).Encode(problem)
}

package main

import (
	"encoding/json"
	"net/http"
)

type errorType struct {
	Text string `json:"errorText"`
}

func respondError(w http.ResponseWriter, code int, errorText string) {
	e := errorType{
		Text: errorText,
	}
	b, _ := json.Marshal(e)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func respondSuccess(w http.ResponseWriter, code int, v interface{}) {
	b, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

package bark

import (
	"encoding/json"
	"net/http"
)

type errorType struct {
	Text string `json:"errorText"`
}

// RespondError is a helper to respond with a standard JSON error response.
func RespondError(w http.ResponseWriter, code int, errorText string) {
	e := errorType{
		Text: errorText,
	}
	b, _ := json.Marshal(e)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

// RespondSuccess is a helper to respond with a successful JSON response.
func RespondSuccess(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

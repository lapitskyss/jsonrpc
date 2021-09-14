package jsonrpc

import (
	"bytes"
	"net/http"

	"github.com/goccy/go-json"
)

type Response struct {
	Version string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id"`
	Result  Result           `json:"result,omitempty"`
	Error   *Error           `json:"error,omitempty"`
}

type Result interface{}

// sendResponse send success JSON response
func sendResponse(w http.ResponseWriter, r []*Response) {
	for i := range r {
		if r[i].Error == nil && r[i].Result == nil {
			r[i].Result = json.RawMessage(nil)
		}
	}

	buf := &bytes.Buffer{}

	if len(r) == 1 {
		if err := json.NewEncoder(buf).Encode(r[0]); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if len(r) > 1 {
		if err := json.NewEncoder(buf).Encode(r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

// sendResponse send single error JSON response
func sendSingleErrorResponse(w http.ResponseWriter, e *Error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	response := &Response{
		Version: Version,
		ID:      nil,
		Error:   e,
	}

	if err := enc.Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

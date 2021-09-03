package jsonrpc

import "github.com/goccy/go-json"

type Request struct {
	Version string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params"`
}

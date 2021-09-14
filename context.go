package jsonrpc

import (
	"net/http"
	"sync"

	"github.com/goccy/go-json"
)

type RequestCtx struct {
	R       *http.Request
	ID      *string
	Version string
	params  json.RawMessage
	mu      sync.RWMutex
	Keys    map[string]interface{}
}

// Params decode request params
// return ErrInvalidParams if can not decode it
func (ctx *RequestCtx) Params(v interface{}) *Error {
	if err := json.Unmarshal(ctx.params, v); err != nil {
		return ErrInvalidParams()
	}
	return nil
}

// Set store a new key/value pair
func (ctx *RequestCtx) Set(key string, value interface{}) {
	ctx.mu.Lock()
	if ctx.Keys == nil {
		ctx.Keys = make(map[string]interface{})
	}

	ctx.Keys[key] = value
	ctx.mu.Unlock()
}

// Get returns the value for the given key,
func (ctx *RequestCtx) Get(key string) (value interface{}, exists bool) {
	ctx.mu.RLock()
	value, exists = ctx.Keys[key]
	ctx.mu.RUnlock()
	return
}

package jsonrpc

import (
	"encoding/json"
	"net/http"
	"sync"
)

type RequestCtx struct {
	R *http.Request

	ID     string
	Params []byte

	mu   sync.RWMutex
	Keys map[string]interface{}
}

// GetParams decode params with standard encoding/json package.
func (ctx *RequestCtx) GetParams(v interface{}) error {
	if err := json.Unmarshal(ctx.Params, v); err != nil {
		return err
	}

	return nil
}

// Result encode json with standard encoding/json package.
func (ctx *RequestCtx) Result(v interface{}) (Result, Error) {
	result, err := json.Marshal(v)
	if err != nil {
		return nil, ErrInternalJSON()
	}

	return result, nil
}

// Set store a new key/value pair.
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

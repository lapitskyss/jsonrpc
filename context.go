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

func (ctx *RequestCtx) Params(v interface{}) *Error {
	if err := json.Unmarshal(ctx.params, v); err != nil {
		return ErrInvalidParams()
	}
	return nil
}

func (ctx *RequestCtx) Set(key string, value interface{}) {
	ctx.mu.Lock()
	if ctx.Keys == nil {
		ctx.Keys = make(map[string]interface{})
	}

	ctx.Keys[key] = value
	ctx.mu.Unlock()
}

func (ctx *RequestCtx) Get(key string) (value interface{}, exists bool) {
	ctx.mu.RLock()
	value, exists = ctx.Keys[key]
	ctx.mu.RUnlock()
	return
}

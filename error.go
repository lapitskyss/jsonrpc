package jsonrpc

import "fmt"

const (
	// ErrorCodeParse Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.
	ErrorCodeParse int = -32700
	// ErrorCodeInvalidRequest The JSON sent is not a valid Request object.
	ErrorCodeInvalidRequest int = -32600
	// ErrorCodeMethodNotFound The method does not exist / is not available.
	ErrorCodeMethodNotFound int = -32601
	// ErrorCodeInvalidParams Invalid method parameter(s).
	ErrorCodeInvalidParams int = -32602
	// ErrorCodeInternal Internal JSON-RPC error.
	ErrorCodeInternal int = -32603
)

type (
	// An Error is a wrapper for a JSON interface value.
	Error struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
)

// Error implements error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("jsonrpc: code: %d, message: %s, data: %+v", e.Code, e.Message, e.Data)
}

// ErrParse returns parse error.
func ErrParse() *Error {
	return &Error{
		Code:    ErrorCodeParse,
		Message: "Parse error",
	}
}

// ErrInvalidRequest returns invalid request error.
func ErrInvalidRequest() *Error {
	return &Error{
		Code:    ErrorCodeInvalidRequest,
		Message: "Invalid Request",
	}
}

// ErrMethodNotFound returns method not found error.
func ErrMethodNotFound() *Error {
	return &Error{
		Code:    ErrorCodeMethodNotFound,
		Message: "Method not found",
	}
}

// ErrInvalidParams returns invalid params error.
func ErrInvalidParams() *Error {
	return &Error{
		Code:    ErrorCodeInvalidParams,
		Message: "Invalid params",
	}
}

// ErrInternal returns internal error.
func ErrInternal() *Error {
	return &Error{
		Code:    ErrorCodeInternal,
		Message: "Internal error",
	}
}

// ErrMaxBatchRequests returns max requests length in batch error.
func ErrMaxBatchRequests() *Error {
	return &Error{
		Code:    ErrorCodeInvalidRequest,
		Message: "Max requests length in batch exceeded",
	}
}

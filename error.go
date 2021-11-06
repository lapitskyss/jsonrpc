package jsonrpc

import (
	"encoding/json"
	"fmt"
)

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
	// ErrorMaxBatchRequests Max requests in batch.
	ErrorMaxBatchRequests int = -32604
)

type Error []byte

// JRPCError is a wrapper for a JSON interface value.
type JRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error implements error interface.
func (e *JRPCError) Error() string {
	return fmt.Sprintf("jsonrpc: code: %d, message: %s, data: %+v", e.Code, e.Message, e.Data)
}

// JSON return json representation of error.
func (e *JRPCError) JSON() []byte {
	result, err := json.Marshal(e)
	if err != nil {
		return ErrInternalJSON()
	}

	return result
}

// ErrParse returns parse error.
func ErrParse() *JRPCError {
	return &JRPCError{
		Code:    ErrorCodeParse,
		Message: "Parse error",
	}
}

// ErrParseJSON return json parse error.
func ErrParseJSON() []byte {
	return []byte(`{"code":-32700,"message":"Parse error"}`)
}

// ErrInvalidRequest returns invalid request error.
func ErrInvalidRequest() *JRPCError {
	return &JRPCError{
		Code:    ErrorCodeInvalidRequest,
		Message: "Invalid Request",
	}
}

// ErrInvalidRequestJSON return json invalid request error.
func ErrInvalidRequestJSON() []byte {
	return []byte(`{"code":-32600,"message":"Invalid Request"}`)
}

// ErrMethodNotFound returns method not found error.
func ErrMethodNotFound() *JRPCError {
	return &JRPCError{
		Code:    ErrorCodeMethodNotFound,
		Message: "Method not found",
	}
}

// ErrMethodNotFoundJSON return json method not found error.
func ErrMethodNotFoundJSON() []byte {
	return []byte(`{"code":-32601,"message":"Method not found"}`)
}

// ErrInvalidParams returns invalid params error.
func ErrInvalidParams() *JRPCError {
	return &JRPCError{
		Code:    ErrorCodeInvalidParams,
		Message: "Invalid params",
	}
}

// ErrInvalidParamsJSON return json invalid params error.
func ErrInvalidParamsJSON() []byte {
	return []byte(`{"code":-32602,"message":"Invalid params"}`)
}

// ErrInternal returns internal error.
func ErrInternal() *JRPCError {
	return &JRPCError{
		Code:    ErrorCodeInternal,
		Message: "Internal error",
	}
}

// ErrInternalJSON return json internal error.
func ErrInternalJSON() []byte {
	return []byte(`{"code":-32603,"message":"Internal error"}`)
}

// ErrMaxBatchRequests returns max requests length in batch error.
func ErrMaxBatchRequests() *JRPCError {
	return &JRPCError{
		Code:    ErrorMaxBatchRequests,
		Message: "Max batch length exceeded",
	}
}

// ErrMaxBatchRequestsJSON return json max requests length in batch error.
func ErrMaxBatchRequestsJSON() []byte {
	return []byte(`{"code":-32604,"message":"Max batch length exceeded"}`)
}

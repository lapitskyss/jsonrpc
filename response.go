package jsonrpc

import (
	"bytes"
	"net/http"
)

type Result []byte

// send result from server.
func send(w http.ResponseWriter, result []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(result)
}

// sendParseError return parse error from server.
func sendParseError(w http.ResponseWriter) {
	send(w, []byte(`{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}`))
}

// sendInvalidRequest return invalid request from server.
func sendInvalidRequest(w http.ResponseWriter) {
	send(w, []byte(`{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null}`))
}

// sendMethodNotFound return method not found from server.
func sendMethodNotFound(w http.ResponseWriter) {
	send(w, []byte(`{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":null}`))
}

// sendInvalidParams return invalid params error from server.
func sendInvalidParams(w http.ResponseWriter) {
	send(w, []byte(`{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":null}`))
}

// sendInternalError return internal error from server.
func sendInternalError(w http.ResponseWriter) {
	send(w, []byte(`{"jsonrpc":"2.0","error":{"code":-32603,"message":"Internal error"},"id":null}`))
}

// sendMaxBatchRequestsError return max batch length exceeded error from server.
func sendMaxBatchRequestsError(w http.ResponseWriter) {
	send(w, []byte(`{"jsonrpc":"2.0","error":{"code":-32604,"message":"Max batch length exceeded"},"id":null}`))
}

// responseMethodNotFound create method not found error response.
func responseMethodNotFound(id []byte) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":`)
	buffer.Write(id)
	buffer.WriteString("}")

	return buffer.Bytes()
}

// responseInvalidRequest create invalid request error response.
func responseInvalidRequest(id []byte) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":`)
	buffer.Write(id)
	buffer.WriteString("}")

	return buffer.Bytes()
}

// responseError create response error with request ID and error.
func responseError(id []byte, err []byte) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`{"jsonrpc":"2.0","error":`)
	buffer.Write(err)
	buffer.WriteString(`,"id":`)
	buffer.Write(id)
	buffer.WriteString("}")

	return buffer.Bytes()
}

// responseResult create response result with request ID and error.
func responseResult(id []byte, result []byte) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`{"jsonrpc":"2.0","result":`)
	buffer.Write(result)
	buffer.WriteString(`,"id":`)
	buffer.Write(id)
	buffer.WriteString("}")

	return buffer.Bytes()
}

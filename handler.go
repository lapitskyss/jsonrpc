package jsonrpc

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/lapitskyss/jsonrpc/jparser"
)

// ServeHTTP process incoming requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), s.options.ContentType) {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	json, err := io.ReadAll(r.Body)
	if err != nil {
		sendInternalError(w)
		return
	}

	if len(json) == 0 {
		sendInvalidRequest(w)
		return
	}

	if err = jparser.ValidateBytes(json); err != nil {
		sendParseError(w)
		return
	}

	if jparser.IsArray(json) {
		batchLen := jparser.ArrayLength(json)
		if batchLen == 0 {
			sendParseError(w)
			return
		}

		if batchLen > s.options.BatchMaxLen {
			sendMaxBatchRequestsError(w)
			return
		}

		respChan := make(chan []byte, batchLen)

		var wg sync.WaitGroup
		wg.Add(batchLen)

		for i := 0; i < batchLen; i++ {
			data := jparser.ArrayElement(json, i)
			go func(data []byte) {
				respChan <- s.handleRequest(r, data)
				wg.Done()
			}(data)
		}

		wg.Wait()
		close(respChan)

		var buffer bytes.Buffer

		buffer.WriteString("[")
		for resp := range respChan {
			buffer.Write(resp)
			buffer.WriteString(",")
		}

		response := buffer.Bytes()
		response[len(response)-1] = ']'

		send(w, response)
		return

	} else {
		send(w, s.handleRequest(r, json))
		return
	}
}

// handleRequest process incoming request single time.
func (s *Server) handleRequest(r *http.Request, json []byte) []byte {
	p := jparser.Parse(json)
	if p.Error() != nil {
		return ErrParseJSON()
	}

	if string(p.Version) != Version {
		return responseInvalidRequest(p.ID)
	}

	method := p.GetMethod()
	if method == "" {
		return responseMethodNotFound(p.ID)
	}

	service := s.GetService(method)
	if service == nil {
		return responseMethodNotFound(p.ID)
	}

	f := service.handler

	for i := len(service.middlewares) - 1; i >= 0; i-- {
		f = service.middlewares[i](f)
	}

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		f = s.middlewares[i](f)
	}

	requestCtx := &RequestCtx{
		R:      r,
		ID:     p.GetId(),
		Params: p.Params,
	}

	result, err := f(requestCtx)
	if err != nil {
		return responseError(p.ID, err)
	}

	return responseResult(p.ID, result)
}

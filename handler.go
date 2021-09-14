package jsonrpc

import (
	"bytes"
	"strings"
	"sync"

	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

// HandleFastHTTP process incoming requests.
func (s *Server) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		return
	}

	if !bytes.HasPrefix(ctx.Request.Header.Peek("Content-Type"), s.options.ContentType) {
		ctx.SetStatusCode(fasthttp.StatusUnsupportedMediaType)
		return
	}

	// parse request
	requests, err := s.parseRequest(ctx)
	if err != nil {
		sendSingleErrorResponse(ctx, err)
		return
	}

	// call global middlewares
	for i := len(s.middlewaresGlobal) - 1; i >= 0; i-- {
		e := s.middlewaresGlobal[i](ctx)
		if e != nil {
			sendSingleErrorResponse(ctx, err)
			return
		}
	}

	// process all requests
	response := s.processRequests(ctx, requests)
	if response == nil {
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	sendResponse(ctx, response)
	return
}

func (s *Server) parseRequest(ctx *fasthttp.RequestCtx) ([]Request, *Error) {
	var requests []Request

	b := ctx.Request.Body()

	if len(b) == 0 {
		return nil, ErrInvalidRequest()
	}

	if b[0] != '[' {
		var req Request
		if err := json.Unmarshal(b, &req); err != nil {
			return nil, ErrParse()
		}
		requests = append(requests, req)
	} else {
		if err := json.Unmarshal(b, &requests); err != nil {
			return nil, ErrParse()
		}
	}

	if len(requests) == 0 {
		return nil, ErrInvalidRequest()
	} else if len(requests) > s.options.BatchMaxLen {
		return nil, ErrMaxBatchRequests()
	}

	return requests, nil
}

func (s *Server) processRequests(ctx *fasthttp.RequestCtx, req []Request) []*Response {
	reqLen := len(req)
	respChan := make(chan *Response, reqLen)

	var wg sync.WaitGroup
	wg.Add(reqLen)

	for i := range req {
		go func(req Request) {
			respChan <- s.processRequest(ctx, req)
			wg.Done()
		}(req[i])
	}

	wg.Wait()
	close(respChan)

	responses := make([]*Response, 0, reqLen)
	for resp := range respChan {
		if resp.ID != nil {
			responses = append(responses, resp)
		}
	}

	if len(responses) == 0 {
		return nil
	}

	return responses
}

func (s *Server) processRequest(ctx *fasthttp.RequestCtx, request Request) *Response {
	if request.Version != Version || request.Method == "" {
		return &Response{
			Version: Version,
			ID:      request.ID,
			Error:   ErrInvalidRequest(),
		}
	}

	method := strings.ToLower(request.Method)
	service, ok := s.services[method]
	if !ok {
		return &Response{
			Version: Version,
			ID:      request.ID,
			Error:   ErrMethodNotFound(),
		}
	}

	f := service.handler

	for i := len(service.middlewares) - 1; i >= 0; i-- {
		f = service.middlewares[i](f)
	}

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		f = s.middlewares[i](f)
	}

	requestCtx := &RequestCtx{
		R:       ctx,
		ID:      getRequestId(request.ID),
		Version: Version,
		params:  request.Params,
	}
	result, err := f(requestCtx)

	return &Response{
		Version: Version,
		ID:      request.ID,
		Result:  result,
		Error:   err,
	}
}

func getRequestId(id *json.RawMessage) *string {
	var result *string
	if id != nil {
		bts := *id
		length := len(bts)
		if length == 0 {
			return nil
		}

		if bts[0] == '"' {
			bts = bts[1 : length-1]
		}

		str := string(bts)
		result = &str
	}

	return result
}

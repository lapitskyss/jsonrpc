package jsonrpc

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/goccy/go-json"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), s.options.ContentType) {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// parse request
	requests, err := s.parseRequest(r)
	if err != nil {
		sendSingleErrorResponse(w, err)
		return
	}

	// call global middlewares
	for i := len(s.middlewaresGlobal) - 1; i >= 0; i-- {
		e := s.middlewaresGlobal[i](r)
		if e != nil {
			sendSingleErrorResponse(w, e)
			return
		}
	}

	// process all requests
	response := s.processRequests(r, requests)
	if response == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	sendResponse(w, response)
	return
}

func (s *Server) parseRequest(r *http.Request) ([]Request, *Error) {
	var requests []Request

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, ErrInvalidRequest()
	}

	err = r.Body.Close()
	if err != nil {
		return nil, ErrInternal()
	}

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

func (s *Server) processRequests(r *http.Request, requests []Request) []Response {
	reqLen := len(requests)
	respChan := make(chan Response, reqLen)

	var wg sync.WaitGroup
	wg.Add(reqLen)

	for _, req := range requests {
		go func(req Request) {
			respChan <- s.processRequest(r, req)
			wg.Done()
		}(req)
	}

	wg.Wait()
	close(respChan)

	responses := make([]Response, 0, reqLen)
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

func (s *Server) processRequest(r *http.Request, request Request) Response {
	if request.Version != Version || request.Method == "" {
		return Response{
			Version: Version,
			ID:      request.ID,
			Error:   ErrInvalidRequest(),
		}
	}

	method := strings.ToLower(request.Method)
	service, ok := s.services[method]
	if !ok {
		return Response{
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
		R:       r,
		ID:      getRequestId(request.ID),
		Version: Version,
		params:  request.Params,
	}
	result, err := f(requestCtx)

	return Response{
		Version: Version,
		ID:      request.ID,
		Result:  result,
		Error:   err,
	}
}

func getRequestId(id *json.RawMessage) *string {
	var result *string
	if id != nil {
		str := strings.Trim(string(*id), `"`)
		result = &str
	}

	return result
}

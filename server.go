package jsonrpc

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	Version            = "2.0"
	defaultBatchMaxLen = 10
	contentTypeJSON    = "application/json"
)

type (
	Handler              func(*RequestCtx) (Result, *Error)
	MiddlewareFunc       func(Handler) Handler
	MiddlewareGlobalFunc func(*http.Request) *Error
)

type Server struct {
	options           Options
	services          map[string]*Service
	middlewares       []MiddlewareFunc
	middlewaresGlobal []MiddlewareGlobalFunc
}

type Service struct {
	handler     Handler
	middlewares []MiddlewareFunc
}

type Options struct {
	BatchMaxLen int
	ContentType string
}

// NewServer create server with provided options
func NewServer(opts Options) *Server {
	if opts.BatchMaxLen == 0 {
		opts.BatchMaxLen = defaultBatchMaxLen
	}

	if opts.ContentType == "" {
		opts.ContentType = contentTypeJSON
	}

	return &Server{
		services: make(map[string]*Service),
		options:  opts,
	}
}

// Register new method
func (s *Server) Register(method string, h Handler) *Service {
	if method == "" {
		panic("can not register service with empty method")
	}

	methodName := strings.ToLower(method)
	if _, ok := s.services[methodName]; ok {
		panic(fmt.Sprintf(`service with name "%s" already exist`, methodName))
	}
	service := &Service{
		handler: h,
	}

	s.services[methodName] = service

	return service
}

// Use appends a middleware handler
func (s *Server) Use(middlewares ...MiddlewareFunc) {
	s.middlewares = append(s.middlewares, middlewares...)
}

// UseGlobal appends a middleware handler. This middleware call ones for all batch requests
func (s *Server) UseGlobal(mGlobal ...MiddlewareGlobalFunc) {
	s.middlewaresGlobal = append(s.middlewaresGlobal, mGlobal...)
}

// Use appends a middleware handler to service
func (service *Service) Use(middlewares ...MiddlewareFunc) {
	service.middlewares = append(service.middlewares, middlewares...)
}

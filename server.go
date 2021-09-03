package jsonrpc

import (
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

func (s *Server) Register(method string, h Handler) *Service {
	service := &Service{
		handler: h,
	}

	s.services[strings.ToLower(method)] = service

	return service
}

func (s *Server) Use(middlewares ...MiddlewareFunc) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) UseGlobal(mGlobal ...MiddlewareGlobalFunc) {
	s.middlewaresGlobal = append(s.middlewaresGlobal, mGlobal...)
}

func (service *Service) Use(middlewares ...MiddlewareFunc) {
	service.middlewares = append(service.middlewares, middlewares...)
}

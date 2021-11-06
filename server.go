package jsonrpc

const (
	Version            = "2.0"
	defaultBatchMaxLen = 10
	contentTypeJSON    = "application/json"
)

type (
	Handler        func(*RequestCtx) (Result, Error)
	MiddlewareFunc func(Handler) Handler
)

type Server struct {
	options     Options
	services    []*Service
	middlewares []MiddlewareFunc
}

type Service struct {
	name        string
	handler     Handler
	middlewares []MiddlewareFunc
}

type Options struct {
	BatchMaxLen int
	ContentType string
}

// NewServer create server with provided options.
func NewServer(opts Options) *Server {
	if opts.BatchMaxLen == 0 {
		opts.BatchMaxLen = defaultBatchMaxLen
	}

	if opts.ContentType == "" {
		opts.ContentType = contentTypeJSON
	}

	return &Server{
		options: opts,
	}
}

// Register new json rpc method.
func (s *Server) Register(method string, h Handler) *Service {
	if method == "" {
		panic("can not register service with empty method")
	}

	service := &Service{
		name:    method,
		handler: h,
	}

	s.services = append(s.services, service)

	return service
}

// GetService get registered service by method name.
func (s *Server) GetService(method string) *Service {
	for _, service := range s.services {
		if service.name == method {
			return service
		}
	}

	return nil
}

// Use appends a middleware handler to server. This middleware call for each service request.
func (s *Server) Use(middlewares ...MiddlewareFunc) {
	s.middlewares = append(s.middlewares, middlewares...)
}

// Use appends a middleware handler to service. This middleware call just for service.
func (service *Service) Use(middlewares ...MiddlewareFunc) {
	service.middlewares = append(service.middlewares, middlewares...)
}

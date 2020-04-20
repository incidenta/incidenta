package server

type Server struct {
	config     Config
	httpServer *HTTPServer
}

func New(c Config) *Server {
	s := &Server{}
	h := NewHTTPServer(c)
	s.config = c
	s.httpServer = h
	return s
}

func (s *Server) Run() error {
	return s.httpServer.Serve()
}

package service

import (
	"net/http"
	_ "test/docs"
)

type Server struct {
	HttpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.HttpServer = &http.Server{Addr: port, Handler: handler}
	return s.HttpServer.ListenAndServe()
}

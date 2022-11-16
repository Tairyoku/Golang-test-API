package test

import (
	"net/http"
	_ "test/docs"
)

type Server struct {
	HttpServer *http.Server
}

// @title          Server API
// @version        0.0.1
// @description    Work with server.
// @termsOfService http://swagger.io/terms/

// @contact.name Tairyoku
// @contact.url  http://www.swagger.io/tairyoku

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     localhost:8080
// @BasePath /
func (s *Server) Run(port string, handler http.Handler) error {
	s.HttpServer = &http.Server{Addr: port, Handler: handler}
	return s.HttpServer.ListenAndServe()
}

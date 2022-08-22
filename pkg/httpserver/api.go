package httpserver

import (
	"fmt"
	"log"
	"net/http"
)

type Interface interface {
	Ping(w http.ResponseWriter, r *http.Request)
}

var _ Interface = (*Server)(nil)

func (s *Server) Ping(w http.ResponseWriter, _ *http.Request) {
	log.Println("Ping request received")
	fmt.Fprintln(w, "Golang HTTP Server")
}

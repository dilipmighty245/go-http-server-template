package httpserver

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewServer(address string, readTimeout, writeTimeout time.Duration) *Server {
	return &Server{Addr: address, ReadTimeout: readTimeout, WriteTimeout: writeTimeout}
}

// server launches a http server
func (s *Server) Start(ctx context.Context) error {
	srv := http.Server{
		Addr:         s.Addr,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		Handler:      s.initHandlers(),
	}

	go func() {
		<-ctx.Done()
		log.Println("attempting graceful shutdown of server")
		srv.SetKeepAlivesEnabled(false)
		closeCtx, closeFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer closeFn()
		_ = srv.Shutdown(closeCtx)
	}()

	return srv.ListenAndServe()
}

func (s *Server) initHandlers() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/ping", s.Ping)
	return r
}

package webserver

import (
	"net/http"
	"api/internal/app/controller/translate"
	"log"
	"context"
)

// WebServer
type Server struct {
}

// NewServer creates new struct.
func NewServer() *Server {
	return &Server{}
}

// ContextHandlerFunc is http.HandlerFunc with context.
type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func newContextHandlerFunc(ctx context.Context, f ContextHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(ctx, w, r)
	}
}

// Run start web server.
func (s Server)Run() {
	ctx := context.Background()

	ctrl := translate.NewController()
	http.HandleFunc("/translate", newContextHandlerFunc(ctx, ctrl.Handle))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
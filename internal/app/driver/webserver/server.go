package webserver

import (
	"api/internal/app/controller/translate"
	"context"
	"log"
	"net/http"
)

// Server represents a web server
type Server struct {
}

// NewServer creates new struct.
func NewServer() *Server {
	return &Server{}
}

// ContextHandlerFunc is http.HandlerFunc with context.
type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func newContextHandlerFunc(ctx context.Context, f ContextHandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.NotFound(w, r)
			return
		}

		f(ctx, w, r)
	}
}

// Run start web server.
func (s Server) Run() {
	ctx := context.Background()

	ctrl := translate.NewController()
	http.HandleFunc("/translate", newContextHandlerFunc(ctx, ctrl.Handle, http.MethodGet))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

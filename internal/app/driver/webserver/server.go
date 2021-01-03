package webserver

import (
	"api/internal/app/controller/translate"
	"api/internal/app/util/appcontext"
	"api/internal/app/util/log"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Server represents a web server
type Server struct {
}

// NewServer creates new struct.
func NewServer() *Server {
	return &Server{}
}

// ContextHandlerFunc is http.HandlerFunc with context.
type ContextHandlerFunc func(appcontext.Context, http.ResponseWriter, *http.Request)

func newContextHandlerFunc(ctx appcontext.Context, f ContextHandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		log.Info(ctx, map[string]interface{}{
			"remoteaddr": r.RemoteAddr,
			"date":       now.Format(time.RFC3339),
			"method":     r.Method,
			"path":       r.URL.Path,
		})

		if r.Method != method {
			http.NotFound(w, r)
			return
		}

		f(ctx, w, r)
	}
}

// Init initialize a server.
func (s Server) Init() {
	log.SetLevel(log.LevelDebug)
}

// Run start web server.
func (s Server) Run() {
	ctx, err := appcontext.NewContext(context.Background(), uuid.New().String())
	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("%v", err))
		return
	}

	ctrl := translate.NewController()
	// TODO: OPTIONS methods
	http.HandleFunc("/translate", newContextHandlerFunc(ctx, ctrl.Handle, http.MethodGet))

	log.Fatal(ctx, http.ListenAndServe(":8080", nil))
}

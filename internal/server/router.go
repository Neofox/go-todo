package server

import (
	"net/http"

	"github.com/a-h/templ"

	"todo/web"
	"todo/web/views"
)

type Middleware func(http.Handler) http.Handler

func applyMiddleware(handler http.Handler, middleware []Middleware) http.Handler {
	for _, m := range middleware {
		handler = m(handler)
	}
	return handler
}

func (s *Server) RegisterRoutes() http.Handler {
	router := http.NewServeMux()

	routes := map[string]struct {
		Handler    http.Handler
		Middleware []Middleware
	}{
		"/": {
			Handler:    http.HandlerFunc(s.handleRoot),
			Middleware: []Middleware{},
		},
		"/todos": {
			Handler:    templ.Handler(views.Test()),
			Middleware: []Middleware{},
		},
		"/helloworld": {
			Handler:    http.HandlerFunc(web.HelloWorld),
			Middleware: []Middleware{},
		},
	}

	for path, handler := range routes {
		router.Handle(path, applyMiddleware(handler.Handler, handler.Middleware))
	}

	fileServer := http.FileServer(http.Dir("static"))
	router.Handle("/static/", http.StripPrefix("/static/", fileServer))

	return router
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World! from handlerfunc"))
}

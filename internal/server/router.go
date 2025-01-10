package server

import (
	"net/http"

	"github.com/a-h/templ"

	"todo/internal/service"
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
			Handler:    templ.Handler(views.Home()),
			Middleware: []Middleware{},
		},
		"/helloworld": {
			Handler:    http.HandlerFunc(web.HelloWorld),
			Middleware: []Middleware{},
		},
		"/todos/": {
			Handler:    http.StripPrefix("/todos", todoHandler(http.NewServeMux())),
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

func todoHandler(router *http.ServeMux) http.Handler {
	controller := web.TodoController{S: *service.NewTodoService()}

	router.HandleFunc("GET /", controller.HandleGet)
	router.HandleFunc("POST /", controller.HandlePost)
	router.HandleFunc("DELETE /{id}", controller.HandleDelete)

	return router
}

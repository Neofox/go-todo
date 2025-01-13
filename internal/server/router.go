package server

import (
	"net/http"

	"github.com/a-h/templ"

	"todo/internal/middleware"
	"todo/internal/service"
	"todo/web/controller"
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
		// templ handler example
		"/": {
			Handler:    templ.Handler(views.Layout("Home", views.Home())),
			Middleware: []Middleware{},
		},
		// resource handler example
		"/todos/": {
			Handler:    http.StripPrefix("/todos", todoHandler(http.NewServeMux())),
			Middleware: []Middleware{middleware.Logger},
		},
	}

	for path, handler := range routes {
		router.Handle(path, applyMiddleware(handler.Handler, handler.Middleware))
	}

	fileServer := http.FileServer(http.Dir("static"))
	router.Handle("/static/", http.StripPrefix("/static/", fileServer))

	return router
}

// todo resource handler
// TODO: maybe move this to a separate file later?
func todoHandler(router *http.ServeMux) http.Handler {
	service := service.NewTodoService()
	controller := controller.TodoController{S: service}

	router.HandleFunc("GET /", controller.HandleGet)
	router.HandleFunc("POST /", controller.HandlePost)
	router.HandleFunc("DELETE /{id}", controller.HandleDelete)

	return router
}

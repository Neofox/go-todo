package server

import (
	"net/http"

	"github.com/a-h/templ"

	"regotth/internal/middleware"
	"regotth/internal/service"
	"regotth/web/controller"
	"regotth/web/view"
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
			Handler:    templ.Handler(view.Layout("Home", view.Home())),
			Middleware: []Middleware{},
		},
		// resource handler example
		"/todos/": {
			Handler:    http.StripPrefix("/todos", todoHandler(http.NewServeMux())),
			Middleware: []Middleware{middleware.Logger},
		},

		"GET /robots.txt": {
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("User-agent: *\nAllow: /"))
			}),
			Middleware: []Middleware{},
		},
	}

	for path, handler := range routes {
		router.Handle(path, applyMiddleware(handler.Handler, handler.Middleware))
	}

	fileServer := http.FileServer(http.Dir("static"))
	router.Handle("/static/", http.StripPrefix("/static", applyMiddleware(fileServer, []Middleware{
		middleware.Cache(middleware.DefaultCacheConfig),
	})))

	return applyMiddleware(router, []Middleware{middleware.Compress})
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

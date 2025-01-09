package web

import (
	"net/http"
	"todo/web/views"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Boby"
	}

	component := views.Hello(name)
	component.Render(r.Context(), w)
}

package web

import (
	"log/slog"
	"net/http"

	"todo/internal/service"
	"todo/web/views"

	e "todo/internal/entity"
)

type TodoController struct {
	S service.TodoService
}

func (c TodoController) HandleGet(w http.ResponseWriter, r *http.Request) {

	todos, err := c.S.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := views.TodoList(todos)
	component.Render(r.Context(), w)
}

func (c TodoController) HandlePost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	// TODO: validate title

	todo, err := c.S.Create(e.Todo{Title: title})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Todo created", "todo", todo)

	component := views.TodoCreated(todo)
	component.Render(r.Context(), w)
}

func (c TodoController) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := c.S.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

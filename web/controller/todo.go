package controller

import (
	"log/slog"
	"net/http"

	e "regotth/internal/entity"
	"regotth/internal/service"
	"regotth/web/view"
)

type TodoController struct {
	S *service.TodoService
}

func (c TodoController) HandleGet(w http.ResponseWriter, r *http.Request) {

	todos, err := c.S.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := view.TodoList(todos)
	view.Layout("Todo List", component).Render(r.Context(), w)
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

	view.TodoCreated(todo).Render(r.Context(), w)
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

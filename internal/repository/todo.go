package repository

import (
	e "todo/internal/entity"

	"github.com/google/uuid"
)

// TodoRepo is a example of a repository implementation for the Todo entity
// in-memory repository implementation
type TodoRepo struct {
	todos map[uuid.UUID]e.Todo
}

// This is the constructor for the TodoRepo struct
// It initializes the repository and returns it
func NewTodoRepo() *TodoRepo {
	return &TodoRepo{
		todos: make(map[uuid.UUID]e.Todo),
	}
}

func (r *TodoRepo) Get(id string) (e.Todo, error) {
	return r.todos[uuid.MustParse(id)], nil
}

func (r *TodoRepo) GetAll() ([]e.Todo, error) {
	var itemList []e.Todo
	for _, item := range r.todos {
		itemList = append(itemList, item)
	}
	return itemList, nil
}

func (r *TodoRepo) Create(t e.Todo) error {
	r.todos[uuid.MustParse(t.ID)] = t
	return nil
}

func (r *TodoRepo) Delete(id string) error {
	delete(r.todos, uuid.MustParse(id))
	return nil
}

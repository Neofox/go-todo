package service

import (
	"todo/internal"
	e "todo/internal/entity"
	r "todo/internal/repository"

	"github.com/google/uuid"
)

// TodoService is a service implementation for the Todo entity
// It is used to define the methods that all services must implement
type TodoService struct {
	repository internal.Repository[e.Todo]
}

func NewTodoService() *TodoService {
	repository := r.NewTodoRepo()
	return &TodoService{repository: repository}
}

func (s *TodoService) Create(todo e.Todo) (e.Todo, error) {

	newTodo := e.Todo{
		ID:        uuid.New().String(),
		Completed: todo.Completed || false,
		Title:     todo.Title,
	}

	return newTodo, s.repository.Create(newTodo)
}

func (s *TodoService) GetAll() ([]e.Todo, error) {
	return s.repository.GetAll()
}

func (s *TodoService) Delete(id string) error {
	return s.repository.Delete(id)
}

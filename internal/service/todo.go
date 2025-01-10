package service

import (
	e "todo/internal/entity"
	r "todo/internal/repository"

	"github.com/google/uuid"
)

type TodoService struct {
	repository r.TodoRepo
}

func NewTodoService() *TodoService {
	return &TodoService{
		repository: *r.NewTodoRepo(),
	}
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

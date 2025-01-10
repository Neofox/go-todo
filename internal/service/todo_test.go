package service

import (
	"testing"
	e "todo/internal/entity"
)

// mockTodoRepo is a mock implementation of the Repository interface
type mockTodoRepo struct {
	todos map[string]e.Todo
}

func newMockTodoRepo() *mockTodoRepo {
	return &mockTodoRepo{
		todos: make(map[string]e.Todo),
	}
}

func (m *mockTodoRepo) GetAll() ([]e.Todo, error) {
	var todos []e.Todo
	for _, todo := range m.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (m *mockTodoRepo) Get(id string) (e.Todo, error) {
	todo, _ := m.todos[id]
	return todo, nil
}

func (m *mockTodoRepo) Create(t e.Todo) error {
	m.todos[t.ID] = t
	return nil
}

func (m *mockTodoRepo) Delete(id string) error {
	delete(m.todos, id)
	return nil
}

// Helper function to create a service with mock repository
func newTestTodoService() (*TodoService, *mockTodoRepo) {
	repo := newMockTodoRepo()
	return &TodoService{repository: repo}, repo
}

func TestTodoService_Create(t *testing.T) {
	svc, repo := newTestTodoService()
	todo := e.Todo{
		Title:     "Test Todo",
		Completed: false,
	}

	created, err := svc.Create(todo)
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	if created.ID == "" {
		t.Error("Expected ID to be set")
	}
	if created.Title != todo.Title {
		t.Errorf("Expected title %q, got %q", todo.Title, created.Title)
	}

	// Verify todo was stored in mock repo
	if len(repo.todos) != 1 {
		t.Error("Expected todo to be stored in repository")
	}
}

func TestTodoService_GetAll(t *testing.T) {
	svc, repo := newTestTodoService()

	// Add a todo directly to mock repo
	todo := e.Todo{ID: "test-id", Title: "Test Todo"}
	repo.todos[todo.ID] = todo

	// Get all todos
	todos, err := svc.GetAll()
	if err != nil {
		t.Fatalf("Failed to get todos: %v", err)
	}

	if len(todos) != 1 {
		t.Error("Expected exactly one todo")
	}
}

func TestTodoService_Delete(t *testing.T) {
	svc, repo := newTestTodoService()

	// Add a todo directly to mock repo
	todo := e.Todo{ID: "test-id", Title: "Test Todo"}
	repo.todos[todo.ID] = todo

	err := svc.Delete(todo.ID)
	if err != nil {
		t.Errorf("Failed to delete todo: %v", err)
	}

	if len(repo.todos) != 0 {
		t.Error("Expected todo to be deleted from repository")
	}
}

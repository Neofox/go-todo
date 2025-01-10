package internal

// Repository is a generic interface for a repository
// It is used to define the methods that all repositories must implement
type Repository[T any] interface {
	GetAll() ([]T, error)
	Get(id string) (T, error)
	Create(t T) error
	Delete(id string) error
}

package internal

type Repository[T any] interface {
	GetAll() ([]T, error)
	Get(id string) (T, error)
	Create(t T) error
	Delete(id string) error
}

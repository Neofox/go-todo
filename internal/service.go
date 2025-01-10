package internal

type Service[T any] struct {
	repository Repository[T]
}

func NewService[T any](repository Repository[T]) *Service[T] {
	return &Service[T]{
		repository: repository,
	}
}

func (s *Service[T]) GetAll() ([]T, error) {
	return s.repository.GetAll()
}

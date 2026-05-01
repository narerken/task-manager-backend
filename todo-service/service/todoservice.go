package service

import (
	"todo-service/models"
	"todo-service/repo"
)

type TodoService struct {
	Repo *repo.TodoRepository
}

func NewTodoService(r *repo.TodoRepository) *TodoService {
	return &TodoService{Repo: r}
}

func (s *TodoService) Create(todo *models.Todo) error {
	return s.Repo.Create(todo)
}

func (s *TodoService) GetAll() ([]models.Todo, error) {
	return s.Repo.GetAll()
}

func (s *TodoService) GetByID(id uint) (*models.Todo, error) {
	return s.Repo.GetByID(id)
}

func (s *TodoService) Update(todo *models.Todo) error {
	return s.Repo.Update(todo)
}

func (s *TodoService) Delete(id uint) error {
	return s.Repo.Delete(id)
}

func (s *TodoService) MarkCompleted(id uint) error {
	todo, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}

	todo.Completed = true
	return s.Repo.Update(todo)
}

func (s *TodoService) GetCompleted() ([]models.Todo, error) {
	return s.Repo.GetCompleted()
}

func (s *TodoService) GetByPriority(priority int) ([]models.Todo, error) {
	return s.Repo.GetByPriority(priority)
}

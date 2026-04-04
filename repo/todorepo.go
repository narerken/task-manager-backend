package repo

import (
	"task-manager/models"

	"gorm.io/gorm"
)

type TodoRepository struct {
	DB *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{DB: db}
}

func (r *TodoRepository) Create(todo *models.Todo) error {
	return r.DB.Create(todo).Error
}

func (r *TodoRepository) GetAll() ([]models.Todo, error) {
	var todos []models.Todo
	err := r.DB.Find(&todos).Error
	return todos, err
}

func (r *TodoRepository) GetByID(id uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.DB.First(&todo, id).Error
	return &todo, err
}

func (r *TodoRepository) Update(todo *models.Todo) error {
	return r.DB.Save(todo).Error
}

func (r *TodoRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Todo{}, id).Error
}

func (r *TodoRepository) GetCompleted() ([]models.Todo, error) {
	var todos []models.Todo
	err := r.DB.Where("completed = ?", true).Find(&todos).Error
	return todos, err
}

func (r *TodoRepository) GetByPriority(priority int) ([]models.Todo, error) {
	var todos []models.Todo
	err := r.DB.Where("priority = ?", priority).Find(&todos).Error
	return todos, err
}

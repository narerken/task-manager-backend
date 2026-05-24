package service

import (
	"context"
	"fmt"
	"task-manager/internal/models"
	"task-manager/internal/repository"
)

type DepartmentServiceInterface interface {
	AddDepartment(ctx context.Context, department *models.Department) (*models.Department, error)
	GetAllDepartments(ctx context.Context) ([]models.Department, error)
	GetDepartmentById(ctx context.Context, id int) (*models.Department, error)
	UpdateDepartment(ctx context.Context, id int, newDepartment models.Department) error
	DeleteDepartment(ctx context.Context, id int) error
}

type DepartmentService struct {
	DepartmentRepo *repository.DepartmentRepository
}

func NewDepartmentService(departmentRepo *repository.DepartmentRepository) *DepartmentService {
	return &DepartmentService{DepartmentRepo: departmentRepo}
}

func (departmentS *DepartmentService) AddDepartment(ctx context.Context, department *models.Department) (*models.Department, error) {
	checkDepartment, err := departmentS.DepartmentRepo.DepartmentExists(ctx, department.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check a department: %w", err)
	}

	if checkDepartment {
		return nil, fmt.Errorf("department already exists")
	}

	err = departmentS.DepartmentRepo.AddDepartment(ctx, department)
	if err != nil {
		return nil, fmt.Errorf("failed to add a department: %w", err)
	}

	return department, nil
}

func (departmentS *DepartmentService) GetAllDepartments(ctx context.Context) ([]models.Department, error) {
	departments, err := departmentS.DepartmentRepo.GetAllDepartments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all departments: %w", err)
	}

	return departments, nil
}

func (departmentS *DepartmentService) GetDepartmentById(ctx context.Context, id int) (*models.Department, error) {
	department, err := departmentS.DepartmentRepo.GetDepartmentById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get a department: %w", err)
	}

	return department, nil
}

func (departmentS *DepartmentService) UpdateDepartment(ctx context.Context, id int, newDepartment *models.Department) error {
	checkDepartment, err := departmentS.DepartmentRepo.DepartmentExists(ctx, newDepartment.Name)
	if err != nil {
		return fmt.Errorf("failed to check a department: %w", err)
	}

	if checkDepartment {
		return fmt.Errorf("department already exists")
	}

	err = departmentS.DepartmentRepo.UpdateDepartment(ctx, id, newDepartment)
	if err != nil {
		return fmt.Errorf("failed to update a department: %w", err)
	}

	return nil
}

func (departmentS *DepartmentService) DeleteDepartment(ctx context.Context, id int) error {
	err := departmentS.DepartmentRepo.DeleteDepartment(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete a department: %w", err)
	}

	return nil
}

package service

import (
	"context"
	"errors"
	"fmt"
	"task-manager/internal/apperrors"
	"task-manager/internal/httpx"
	"task-manager/internal/models"
	"task-manager/internal/models/inputs"
	"task-manager/internal/repository"
)

type AttendanceServiceInterface interface {
	AddAttendance(ctx context.Context, attendance *models.Attendance) (*models.Attendance, error)
	GetAllAttendances(ctx context.Context) ([]models.Attendance, error)
	GetAttendanceById(ctx context.Context, id int) (*models.Attendance, error)
	UpdateAttendance(ctx context.Context, id int, in *inputs.UpdateAttendanceInput) error
	DeleteAttendance(ctx context.Context, id int) error
}

type AttendanceService struct {
	AttendanceRepo *repository.AttendanceRepository
	UserS          *UserService
}

func NewAttendanceService(attendanceR *repository.AttendanceRepository, userS *UserService) *AttendanceService {
	return &AttendanceService{
		AttendanceRepo: attendanceR,
		UserS:          userS,
	}
}

func (attendanceS *AttendanceService) AddAttendance(ctx context.Context, in *inputs.AddAttendanceInput) (*models.Attendance, error) {
	creatorUser, err := attendanceS.UserS.GetUserById(ctx, in.Creator)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, httpx.NotFound("user")
		}
	}

	user, err := attendanceS.UserS.GetUserById(ctx, in.UserId)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, httpx.NotFound("user")
		}

		return nil, fmt.Errorf("failed to get creator user: %w", err)
	}

	if *creatorUser.DepartmentId != *user.DepartmentId {
		return nil, fmt.Errorf("creator user and user must be in the same department")
	}

	attendance := models.Attendance{
		Date:         in.Date,
		UserId:       in.UserId,
		DepartmentId: user.DepartmentId,
		Comment:      in.Comment,
	}

	err = attendanceS.AttendanceRepo.AddAttendance(ctx, &attendance)
	if err != nil {
		return nil, fmt.Errorf("failed to add an attendance: %w", err)
	}

	return &attendance, nil
}

func (attendanceS *AttendanceService) GetAllAttendances(ctx context.Context) ([]models.Attendance, error) {
	attendances, err := attendanceS.AttendanceRepo.GetAllAttendances(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all attendances: %w", err)
	}

	return attendances, nil
}

func (attendanceS *AttendanceService) GetAttendanceById(ctx context.Context, id int) (*models.Attendance, error) {
	attendance, err := attendanceS.AttendanceRepo.GetAttendanceById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get an attendance: %w", err)
	}

	return attendance, nil
}

func (attendanceS *AttendanceService) UpdateAttendance(ctx context.Context, id int, in *inputs.UpdateAttendanceInput) error {
	err := attendanceS.AttendanceRepo.UpdateAttendance(ctx, id, in)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrNotFound
		}

		return fmt.Errorf("failed to update an attendance: %w", err)
	}

	return nil
}

func (attendanceS *AttendanceService) DeleteAttendance(ctx context.Context, id int) error {
	err := attendanceS.AttendanceRepo.DeleteAttendance(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete an attendance: %w", err)
	}

	return nil
}

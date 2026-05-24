package repository

import (
	"context"
	"errors"
	"fmt"
	"task-manager/internal/apperrors"
	"task-manager/internal/models"
	"task-manager/internal/models/inputs"

	"gorm.io/gorm"
)

type AttendanceRepositoryInterface interface {
	AddAttendance(ctx context.Context, attendance *models.Attendance) error
	GetAttendanceById(ctx context.Context, id int) (*models.Attendance, error)
	GetAllAttendances(ctx context.Context) ([]models.Attendance, error)
	UpdateAttendance(ctx context.Context, id int, in *inputs.UpdateAttendanceInput) error
	DeleteAttendance(ctx context.Context, id int) error
}

type AttendanceRepository struct {
	DB *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{DB: db}
}

func (attendanceRepo *AttendanceRepository) AddAttendance(ctx context.Context, attendance *models.Attendance) error {
	err := attendanceRepo.DB.WithContext(ctx).Table("attendance").Create(attendance).Error
	if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	return nil
}

func (attendanceRepo *AttendanceRepository) GetAttendanceById(ctx context.Context, id int) (*models.Attendance, error) {
	var attendance models.Attendance
	err := attendanceRepo.DB.WithContext(ctx).
		Table("attendance").
		Select("id", "date", "user_id", "department_id", "status", "comment", "marked_by", "created_at", "updated_at").
		Where("id = ? AND deleted_at IS NULL", id).
		Take(&attendance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &attendance, nil
}

func (attendanceRepo *AttendanceRepository) GetAllAttendances(ctx context.Context) ([]models.Attendance, error) {
	var attendances []models.Attendance
	err := attendanceRepo.DB.WithContext(ctx).
		Table("attendance").
		Select("id", "date", "user_id", "department_id", "status", "comment", "marked_by", "created_at", "updated_at").
		Where("deleted_at IS NULL").
		Find(&attendances).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all attendances: %w", err)
	}

	return attendances, nil
}

func (attendanceRepo *AttendanceRepository) UpdateAttendance(ctx context.Context, id int, in *inputs.UpdateAttendanceInput) error {
	updates := map[string]any{}
	if in.Status != nil {
		updates["status"] = *in.Status
		if in.MarkedBy != nil {
			updates["marked_by"] = *in.MarkedBy
		}
	}
	if in.Comment != nil {
		updates["comment"] = *in.Comment
	}
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	res := attendanceRepo.DB.WithContext(ctx).
		Table("attendance").
		Where("id = ?", id).
		UpdateColumns(updates)
	if res.Error != nil {
		return fmt.Errorf("failed to update attendance: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (attendanceRepo *AttendanceRepository) DeleteAttendance(ctx context.Context, id int) error {
	res := attendanceRepo.DB.WithContext(ctx).
		Table("attendance").
		Where("id = ? AND deleted_at IS NULL", id).
		UpdateColumn("deleted_at", gorm.Expr("now()"))
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	if res.Error != nil {
		return fmt.Errorf("failed to delete an attendance: %w", res.Error)
	}

	return nil
}

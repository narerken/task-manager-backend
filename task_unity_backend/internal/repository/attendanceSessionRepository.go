package repository

import (
	"context"
	"errors"
	"fmt"
	"task-manager/internal/apperrors"
	"task-manager/internal/models"
	"time"

	"gorm.io/gorm"
)

type AttendanceSessionRepository struct {
	DB *gorm.DB
}

func NewAttendanceSessionRepository(db *gorm.DB) *AttendanceSessionRepository {
	return &AttendanceSessionRepository{DB: db}
}

func (repo *AttendanceSessionRepository) CreateSession(ctx context.Context, session *models.AttendanceSession) error {
	if err := repo.DB.WithContext(ctx).Table("attendance_sessions").Create(session).Error; err != nil {
		return fmt.Errorf("failed to create attendance session: %w", err)
	}
	return nil
}

func (repo *AttendanceSessionRepository) GetSessionById(ctx context.Context, id int) (*models.AttendanceSession, error) {
	var session models.AttendanceSession
	err := repo.DB.WithContext(ctx).
		Table("attendance_sessions").
		Select("id", "department_id", "date", "state", "created_by", "updated_by", "created_at", "updated_at", "deleted_at").
		Where("id = ? AND deleted_at IS NULL", id).
		Take(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get attendance session: %w", err)
	}
	return &session, nil
}

func (repo *AttendanceSessionRepository) FindSessionByDepartmentAndDate(ctx context.Context, departmentId int, date time.Time) (*models.AttendanceSession, error) {
	var session models.AttendanceSession
	err := repo.DB.WithContext(ctx).
		Table("attendance_sessions").
		Select("id", "department_id", "date", "state", "created_by", "updated_by", "created_at", "updated_at", "deleted_at").
		Where("department_id = ? AND date = ? AND deleted_at IS NULL", departmentId, date).
		Take(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find attendance session: %w", err)
	}
	return &session, nil
}

func (repo *AttendanceSessionRepository) ListSessions(ctx context.Context, departmentId *int, dateFrom, dateTo *time.Time) ([]models.AttendanceSession, error) {
	query := repo.DB.WithContext(ctx).
		Table("attendance_sessions").
		Select("id", "department_id", "date", "state", "created_by", "updated_by", "created_at", "updated_at", "deleted_at").
		Where("deleted_at IS NULL")

	if departmentId != nil {
		query = query.Where("department_id = ?", *departmentId)
	}
	if dateFrom != nil {
		query = query.Where("date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("date <= ?", *dateTo)
	}

	var sessions []models.AttendanceSession
	if err := query.Order("date DESC, id DESC").Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to list attendance sessions: %w", err)
	}

	return sessions, nil
}

func (repo *AttendanceSessionRepository) UpdateSessionState(ctx context.Context, sessionId int, state string, actorId int) error {
	res := repo.DB.WithContext(ctx).
		Table("attendance_sessions").
		Where("id = ? AND deleted_at IS NULL", sessionId).
		UpdateColumns(map[string]any{
			"state":      state,
			"updated_by": actorId,
			"updated_at": gorm.Expr("now()"),
		})
	if res.Error != nil {
		return fmt.Errorf("failed to update attendance session state: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (repo *AttendanceSessionRepository) GetSessionEntriesView(ctx context.Context, sessionId int, departmentId int) ([]models.AttendanceSessionEntryView, error) {
	type row struct {
		StudentId   int
		StudentName string
		Status      *string
		Comment     *string
		MarkedBy    *int
	}

	var rows []row
	err := repo.DB.WithContext(ctx).
		Table("users").
		Select("users.id as student_id", "users.full_name as student_name", "attendance_entries.status", "attendance_entries.comment", "attendance_entries.marked_by").
		Joins("LEFT JOIN attendance_entries ON attendance_entries.student_id = users.id AND attendance_entries.session_id = ? AND attendance_entries.deleted_at IS NULL", sessionId).
		Where("users.department_id = ? AND users.deleted_at IS NULL", departmentId).
		Order("users.full_name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get attendance session entries: %w", err)
	}

	entries := make([]models.AttendanceSessionEntryView, 0, len(rows))
	for _, r := range rows {
		entries = append(entries, models.AttendanceSessionEntryView{
			StudentId:   r.StudentId,
			StudentName: r.StudentName,
			Status:      r.Status,
			Comment:     r.Comment,
			MarkedBy:    r.MarkedBy,
		})
	}

	return entries, nil
}

func (repo *AttendanceSessionRepository) UpsertEntries(ctx context.Context, sessionId int, actorId int, entries []models.AttendanceEntry) error {
	return repo.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, entry := range entries {
			var existing models.AttendanceEntry
			err := tx.Table("attendance_entries").
				Select("id", "session_id", "student_id", "status", "comment", "marked_by", "created_at", "updated_at", "deleted_at").
				Where("session_id = ? AND student_id = ? AND deleted_at IS NULL", sessionId, entry.StudentId).
				Take(&existing).Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					newEntry := models.AttendanceEntry{
						SessionId: sessionId,
						StudentId: entry.StudentId,
						Status:    entry.Status,
						Comment:   entry.Comment,
						MarkedBy:  &actorId,
					}
					if createErr := tx.Table("attendance_entries").Create(&newEntry).Error; createErr != nil {
						return fmt.Errorf("failed to create attendance entry: %w", createErr)
					}
					continue
				}
				return fmt.Errorf("failed to read attendance entry: %w", err)
			}

			res := tx.Table("attendance_entries").
				Where("id = ? AND deleted_at IS NULL", existing.Id).
				UpdateColumns(map[string]any{
					"status":     entry.Status,
					"comment":    entry.Comment,
					"marked_by":  actorId,
					"updated_at": gorm.Expr("now()"),
				})
			if res.Error != nil {
				return fmt.Errorf("failed to update attendance entry: %w", res.Error)
			}
		}

		res := tx.Table("attendance_sessions").
			Where("id = ? AND deleted_at IS NULL", sessionId).
			UpdateColumns(map[string]any{
				"updated_by": actorId,
				"updated_at": gorm.Expr("now()"),
			})
		if res.Error != nil {
			return fmt.Errorf("failed to bump attendance session timestamp: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return apperrors.ErrNotFound
		}

		return nil
	})
}

package repository

import (
	"attendance_session_service/internal/models"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

type AttendanceSessionRepository struct {
	db *gorm.DB
}

func NewAttendanceSessionRepository(db *gorm.DB) *AttendanceSessionRepository {
	return &AttendanceSessionRepository{db: db}
}

func (r *AttendanceSessionRepository) CreateSession(ctx context.Context, session *models.AttendanceSession) error {
	if err := r.db.WithContext(ctx).Table("attendance_sessions").Create(session).Error; err != nil {
		return fmt.Errorf("failed to create attendance session: %w", err)
	}
	return nil
}

func (r *AttendanceSessionRepository) GetSessionByID(ctx context.Context, id int) (*models.AttendanceSession, error) {
	var session models.AttendanceSession
	err := r.db.WithContext(ctx).
		Table("attendance_sessions").
		Where("id = ? AND deleted_at IS NULL", id).
		Take(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get attendance session: %w", err)
	}
	return &session, nil
}

func (r *AttendanceSessionRepository) FindSessionByDepartmentAndDate(ctx context.Context, departmentID int, date time.Time) (*models.AttendanceSession, error) {
	var session models.AttendanceSession
	err := r.db.WithContext(ctx).
		Table("attendance_sessions").
		Where("department_id = ? AND date = ? AND deleted_at IS NULL", departmentID, date).
		Take(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find attendance session: %w", err)
	}
	return &session, nil
}

func (r *AttendanceSessionRepository) ListSessions(ctx context.Context, departmentID *int, dateFrom, dateTo *time.Time) ([]models.AttendanceSession, error) {
	query := r.db.WithContext(ctx).Table("attendance_sessions").Where("deleted_at IS NULL")
	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
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

func (r *AttendanceSessionRepository) UpdateSessionState(ctx context.Context, sessionID int, state string, actorID int) error {
	res := r.db.WithContext(ctx).
		Table("attendance_sessions").
		Where("id = ? AND deleted_at IS NULL", sessionID).
		UpdateColumns(map[string]any{
			"state":      state,
			"updated_by": actorID,
			"updated_at": gorm.Expr("now()"),
		})
	if res.Error != nil {
		return fmt.Errorf("failed to update attendance session state: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *AttendanceSessionRepository) GetEntriesBySessionID(ctx context.Context, sessionID int) ([]models.AttendanceEntry, error) {
	var entries []models.AttendanceEntry
	if err := r.db.WithContext(ctx).
		Table("attendance_entries").
		Where("session_id = ? AND deleted_at IS NULL", sessionID).
		Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("failed to get attendance entries: %w", err)
	}
	return entries, nil
}

func (r *AttendanceSessionRepository) UpsertEntries(ctx context.Context, sessionID int, actorID int, entries []models.AttendanceEntry) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, entry := range entries {
			var existing models.AttendanceEntry
			err := tx.Table("attendance_entries").
				Where("session_id = ? AND student_id = ? AND deleted_at IS NULL", sessionID, entry.StudentId).
				Take(&existing).Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					entry.SessionId = sessionID
					entry.MarkedBy = &actorID
					if createErr := tx.Table("attendance_entries").Create(&entry).Error; createErr != nil {
						return fmt.Errorf("failed to create attendance entry: %w", createErr)
					}
					continue
				}
				return fmt.Errorf("failed to read attendance entry: %w", err)
			}

			res := tx.Table("attendance_entries").
				Where("id = ?", existing.Id).
				UpdateColumns(map[string]any{
					"status":     entry.Status,
					"comment":    entry.Comment,
					"marked_by":  actorID,
					"updated_at": gorm.Expr("now()"),
				})
			if res.Error != nil {
				return fmt.Errorf("failed to update attendance entry: %w", res.Error)
			}
		}

		res := tx.Table("attendance_sessions").
			Where("id = ? AND deleted_at IS NULL", sessionID).
			UpdateColumns(map[string]any{
				"updated_by": actorID,
				"updated_at": gorm.Expr("now()"),
			})
		if res.Error != nil {
			return fmt.Errorf("failed to bump session timestamp: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return ErrNotFound
		}

		return nil
	})
}

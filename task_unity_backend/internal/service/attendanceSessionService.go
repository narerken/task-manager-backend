package service

import (
	"context"
	"errors"
	"fmt"
	"task-manager/internal/apperrors"
	"task-manager/internal/models"
	"task-manager/internal/models/inputs"
	"task-manager/internal/repository"
	"time"
)

type AttendanceSessionService struct {
	SessionRepo    *repository.AttendanceSessionRepository
	UserRepo       *repository.UserRepository
	DepartmentRepo *repository.DepartmentRepository
}

func NewAttendanceSessionService(sessionRepo *repository.AttendanceSessionRepository, userRepo *repository.UserRepository, departmentRepo *repository.DepartmentRepository) *AttendanceSessionService {
	return &AttendanceSessionService{
		SessionRepo:    sessionRepo,
		UserRepo:       userRepo,
		DepartmentRepo: departmentRepo,
	}
}

func (s *AttendanceSessionService) CreateSession(ctx context.Context, currentUserId int, role string, in inputs.CreateAttendanceSessionInput) (*models.AttendanceSession, error) {
	sessionDate, err := parseAttendanceDate(in.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date: %w", err)
	}

	currentUser, err := s.UserRepo.GetUserById(ctx, currentUserId)
	if err != nil {
		return nil, fmt.Errorf("failed to load current user: %w", err)
	}

	if err := s.ensureDepartmentAccess(currentUser, role, in.DepartmentId); err != nil {
		return nil, err
	}

	if _, err := s.DepartmentRepo.GetDepartmentById(ctx, in.DepartmentId); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to validate department: %w", err)
	}

	existing, err := s.SessionRepo.FindSessionByDepartmentAndDate(ctx, in.DepartmentId, sessionDate)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, apperrors.ErrNotFound) {
		return nil, fmt.Errorf("failed to check existing attendance session: %w", err)
	}

	session := models.AttendanceSession{
		DepartmentId: in.DepartmentId,
		Date:         sessionDate,
		State:        "draft",
		CreatedBy:    intPtr(currentUserId),
		UpdatedBy:    intPtr(currentUserId),
	}
	if err := s.SessionRepo.CreateSession(ctx, &session); err != nil {
		return nil, fmt.Errorf("failed to create attendance session: %w", err)
	}
	return &session, nil
}

func (s *AttendanceSessionService) GetSession(ctx context.Context, currentUserId int, role string, sessionId int) (*models.AttendanceSessionDetails, error) {
	session, err := s.SessionRepo.GetSessionById(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	currentUser, err := s.UserRepo.GetUserById(ctx, currentUserId)
	if err != nil {
		return nil, fmt.Errorf("failed to load current user: %w", err)
	}

	if err := s.ensureDepartmentAccess(currentUser, role, session.DepartmentId); err != nil {
		return nil, err
	}

	entries, err := s.SessionRepo.GetSessionEntriesView(ctx, session.Id, session.DepartmentId)
	if err != nil {
		return nil, fmt.Errorf("failed to get session entries: %w", err)
	}

	return &models.AttendanceSessionDetails{
		Id:           session.Id,
		DepartmentId: session.DepartmentId,
		Date:         session.Date.Format("2006-01-02"),
		State:        session.State,
		CreatedBy:    session.CreatedBy,
		UpdatedBy:    session.UpdatedBy,
		Entries:      entries,
	}, nil
}

func (s *AttendanceSessionService) ListSessions(ctx context.Context, currentUserId int, role string, departmentId *int, dateFrom, dateTo *string) ([]models.AttendanceSession, error) {
	currentUser, err := s.UserRepo.GetUserById(ctx, currentUserId)
	if err != nil {
		return nil, fmt.Errorf("failed to load current user: %w", err)
	}

	var from, to *time.Time
	if dateFrom != nil && *dateFrom != "" {
		parsed, parseErr := parseAttendanceDate(*dateFrom)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid date_from: %w", parseErr)
		}
		from = &parsed
	}
	if dateTo != nil && *dateTo != "" {
		parsed, parseErr := parseAttendanceDate(*dateTo)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid date_to: %w", parseErr)
		}
		to = &parsed
	}

	filterDepartment := departmentId
	if role != "admin" {
		if currentUser.DepartmentId == nil {
			return nil, fmt.Errorf("current user department is not set")
		}
		if departmentId != nil && *departmentId != *currentUser.DepartmentId {
			return nil, apperrors.ErrUnauthorized
		}
		filterDepartment = currentUser.DepartmentId
	}

	sessions, err := s.SessionRepo.ListSessions(ctx, filterDepartment, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to list attendance sessions: %w", err)
	}
	return sessions, nil
}

func (s *AttendanceSessionService) BulkUpsertEntries(ctx context.Context, currentUserId int, role string, sessionId int, in inputs.BulkUpsertAttendanceEntriesInput) error {
	session, err := s.SessionRepo.GetSessionById(ctx, sessionId)
	if err != nil {
		return err
	}

	currentUser, err := s.UserRepo.GetUserById(ctx, currentUserId)
	if err != nil {
		return fmt.Errorf("failed to load current user: %w", err)
	}

	if err := s.ensureDepartmentAccess(currentUser, role, session.DepartmentId); err != nil {
		return err
	}

	if session.State == "published" {
		return fmt.Errorf("attendance session is already published")
	}

	entries := make([]models.AttendanceEntry, 0, len(in.Entries))
	for _, item := range in.Entries {
		student, getErr := s.UserRepo.GetUserById(ctx, item.StudentId)
		if getErr != nil {
			if errors.Is(getErr, apperrors.ErrNotFound) {
				return apperrors.ErrNotFound
			}
			return fmt.Errorf("failed to load student: %w", getErr)
		}
		if student.DepartmentId == nil || *student.DepartmentId != session.DepartmentId {
			return fmt.Errorf("student %d does not belong to session department", item.StudentId)
		}

		entries = append(entries, models.AttendanceEntry{
			SessionId: sessionId,
			StudentId: item.StudentId,
			Status:    item.Status,
			Comment:   item.Comment,
		})
	}

	if err := s.SessionRepo.UpsertEntries(ctx, sessionId, currentUserId, entries); err != nil {
		return fmt.Errorf("failed to upsert attendance entries: %w", err)
	}

	return nil
}

func (s *AttendanceSessionService) PublishSession(ctx context.Context, currentUserId int, role string, sessionId int) error {
	session, err := s.SessionRepo.GetSessionById(ctx, sessionId)
	if err != nil {
		return err
	}

	currentUser, err := s.UserRepo.GetUserById(ctx, currentUserId)
	if err != nil {
		return fmt.Errorf("failed to load current user: %w", err)
	}

	if err := s.ensureDepartmentAccess(currentUser, role, session.DepartmentId); err != nil {
		return err
	}

	if err := s.SessionRepo.UpdateSessionState(ctx, sessionId, "published", currentUserId); err != nil {
		return fmt.Errorf("failed to publish attendance session: %w", err)
	}
	return nil
}

func (s *AttendanceSessionService) ensureDepartmentAccess(currentUser *models.User, role string, departmentId int) error {
	if role == "admin" {
		return nil
	}
	if role != "manager" {
		return apperrors.ErrUnauthorized
	}
	if currentUser.DepartmentId == nil {
		return fmt.Errorf("current user department is not set")
	}
	if *currentUser.DepartmentId != departmentId {
		return apperrors.ErrUnauthorized
	}
	return nil
}

func parseAttendanceDate(value string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}

func intPtr(v int) *int {
	return &v
}

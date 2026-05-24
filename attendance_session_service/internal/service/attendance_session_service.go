package service

import (
	"attendance_session_service/internal/client"
	"attendance_session_service/internal/models"
	"attendance_session_service/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"
)

type AttendanceSessionService struct {
	repo       *repository.AttendanceSessionRepository
	mainClient *client.MainServiceClient
}

func NewAttendanceSessionService(repo *repository.AttendanceSessionRepository, mainClient *client.MainServiceClient) *AttendanceSessionService {
	return &AttendanceSessionService{repo: repo, mainClient: mainClient}
}

func (s *AttendanceSessionService) CreateSession(ctx context.Context, currentUserID int, role string, in models.CreateAttendanceSessionInput) (*models.AttendanceSession, error) {
	sessionDate, err := parseAttendanceDate(in.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date: %w", err)
	}

	currentUser, err := s.mainClient.GetUserByID(ctx, currentUserID)
	if err != nil {
		return nil, mapRemoteError("failed to load current user", err)
	}

	if err := s.ensureDepartmentAccess(currentUser, role, in.DepartmentId); err != nil {
		return nil, err
	}

	if _, err := s.mainClient.GetDepartmentByID(ctx, in.DepartmentId); err != nil {
		return nil, mapRemoteError("failed to validate department", err)
	}

	existing, err := s.repo.FindSessionByDepartmentAndDate(ctx, in.DepartmentId, sessionDate)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("failed to check existing attendance session: %w", err)
	}

	session := models.AttendanceSession{
		DepartmentId: in.DepartmentId,
		Date:         sessionDate,
		State:        "draft",
		CreatedBy:    intPtr(currentUserID),
		UpdatedBy:    intPtr(currentUserID),
	}
	if err := s.repo.CreateSession(ctx, &session); err != nil {
		return nil, fmt.Errorf("failed to create attendance session: %w", err)
	}
	return &session, nil
}

func (s *AttendanceSessionService) GetSession(ctx context.Context, currentUserID int, role string, sessionID int) (*models.AttendanceSessionDetails, error) {
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	currentUser, err := s.mainClient.GetUserByID(ctx, currentUserID)
	if err != nil {
		return nil, mapRemoteError("failed to load current user", err)
	}
	if err := s.ensureDepartmentAccess(currentUser, role, session.DepartmentId); err != nil {
		return nil, err
	}

	users, err := s.mainClient.GetDepartmentUsers(ctx, session.DepartmentId)
	if err != nil {
		return nil, mapRemoteError("failed to load department users", err)
	}

	entryRows, err := s.repo.GetEntriesBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session entries: %w", err)
	}

	entriesByStudentID := make(map[int]models.AttendanceEntry, len(entryRows))
	for _, entry := range entryRows {
		entriesByStudentID[entry.StudentId] = entry
	}

	entries := make([]models.AttendanceSessionEntryView, 0, len(users))
	for _, user := range users {
		entry, ok := entriesByStudentID[user.Id]
		if ok {
			entries = append(entries, models.AttendanceSessionEntryView{
				StudentId:   user.Id,
				StudentName: user.FullName,
				Status:      &entry.Status,
				Comment:     entry.Comment,
				MarkedBy:    entry.MarkedBy,
			})
			continue
		}

		entries = append(entries, models.AttendanceSessionEntryView{
			StudentId:   user.Id,
			StudentName: user.FullName,
		})
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

func (s *AttendanceSessionService) ListSessions(ctx context.Context, currentUserID int, role string, departmentID *int, dateFrom, dateTo *string) ([]models.AttendanceSession, error) {
	currentUser, err := s.mainClient.GetUserByID(ctx, currentUserID)
	if err != nil {
		return nil, mapRemoteError("failed to load current user", err)
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

	filterDepartment := departmentID
	if role != "admin" {
		if currentUser.DepartmentId == nil {
			return nil, fmt.Errorf("current user department is not set")
		}
		if departmentID != nil && *departmentID != *currentUser.DepartmentId {
			return nil, ErrUnauthorized
		}
		filterDepartment = currentUser.DepartmentId
	}

	return s.repo.ListSessions(ctx, filterDepartment, from, to)
}

func (s *AttendanceSessionService) BulkUpsertEntries(ctx context.Context, currentUserID int, role string, sessionID int, in models.BulkUpsertAttendanceEntriesInput) error {
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}

	currentUser, err := s.mainClient.GetUserByID(ctx, currentUserID)
	if err != nil {
		return mapRemoteError("failed to load current user", err)
	}
	if err := s.ensureDepartmentAccess(currentUser, role, session.DepartmentId); err != nil {
		return err
	}
	if session.State == "published" {
		return fmt.Errorf("attendance session is already published")
	}

	users, err := s.mainClient.GetDepartmentUsers(ctx, session.DepartmentId)
	if err != nil {
		return mapRemoteError("failed to load department users", err)
	}
	userByID := make(map[int]models.User, len(users))
	for _, user := range users {
		userByID[user.Id] = user
	}

	entries := make([]models.AttendanceEntry, 0, len(in.Entries))
	for _, item := range in.Entries {
		if _, ok := userByID[item.StudentId]; !ok {
			return fmt.Errorf("student %d does not belong to session department", item.StudentId)
		}
		entries = append(entries, models.AttendanceEntry{
			SessionId: sessionID,
			StudentId: item.StudentId,
			Status:    item.Status,
			Comment:   item.Comment,
		})
	}

	return s.repo.UpsertEntries(ctx, sessionID, currentUserID, entries)
}

func (s *AttendanceSessionService) PublishSession(ctx context.Context, currentUserID int, role string, sessionID int) error {
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}

	currentUser, err := s.mainClient.GetUserByID(ctx, currentUserID)
	if err != nil {
		return mapRemoteError("failed to load current user", err)
	}
	if err := s.ensureDepartmentAccess(currentUser, role, session.DepartmentId); err != nil {
		return err
	}

	return s.repo.UpdateSessionState(ctx, sessionID, "published", currentUserID)
}

var ErrUnauthorized = errors.New("unauthorized")

func (s *AttendanceSessionService) ensureDepartmentAccess(currentUser *models.User, role string, departmentID int) error {
	if role == "admin" {
		return nil
	}
	if role != "manager" {
		return ErrUnauthorized
	}
	if currentUser.DepartmentId == nil {
		return fmt.Errorf("current user department is not set")
	}
	if *currentUser.DepartmentId != departmentID {
		return ErrUnauthorized
	}
	return nil
}

func parseAttendanceDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

func intPtr(v int) *int {
	return &v
}

func mapRemoteError(message string, err error) error {
	if errors.Is(err, client.ErrRemoteNotFound) {
		return repository.ErrNotFound
	}
	return fmt.Errorf("%s: %w", message, err)
}

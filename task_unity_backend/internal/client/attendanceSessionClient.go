package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type AttendanceSessionClient struct {
	baseURL      string
	serviceToken string
	client       *resty.Client
}

func NewAttendanceSessionClient(baseURL, serviceToken string, timeoutSeconds int) *AttendanceSessionClient {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 10
	}

	return &AttendanceSessionClient{
		baseURL:      strings.TrimRight(baseURL, "/"),
		serviceToken: serviceToken,
		client: resty.New().
			SetTimeout(time.Duration(timeoutSeconds) * time.Second),
	}
}

func (c *AttendanceSessionClient) CreateSession(ctx context.Context, userID int, role string, body []byte) (*resty.Response, error) {
	return c.newRequest(ctx, userID, role).
		SetBody(body).
		Post(c.baseURL + "/attendance/sessions")
}

func (c *AttendanceSessionClient) GetSession(ctx context.Context, userID int, role string, sessionID string) (*resty.Response, error) {
	return c.newRequest(ctx, userID, role).
		Get(c.baseURL + "/attendance/sessions/" + sessionID)
}

func (c *AttendanceSessionClient) ListSessions(ctx context.Context, userID int, role string, query url.Values) (*resty.Response, error) {
	return c.newRequest(ctx, userID, role).
		SetQueryParamsFromValues(query).
		Get(c.baseURL + "/attendance/sessions")
}

func (c *AttendanceSessionClient) BulkUpsertEntries(ctx context.Context, userID int, role string, sessionID string, body []byte) (*resty.Response, error) {
	return c.newRequest(ctx, userID, role).
		SetBody(body).
		Patch(c.baseURL + "/attendance/sessions/" + sessionID + "/entries")
}

func (c *AttendanceSessionClient) PublishSession(ctx context.Context, userID int, role string, sessionID string) (*resty.Response, error) {
	return c.newRequest(ctx, userID, role).
		Patch(c.baseURL + "/attendance/sessions/" + sessionID + "/publish")
}

func (c *AttendanceSessionClient) newRequest(ctx context.Context, userID int, role string) *resty.Request {
	return c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Service-Token", c.serviceToken).
		SetHeader("X-User-Id", strconv.Itoa(userID)).
		SetHeader("X-User-Role", role)
}

func ProxyError(err error) error {
	return fmt.Errorf("attendance session service request failed: %w", err)
}

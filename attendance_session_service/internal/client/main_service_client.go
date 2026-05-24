package client

import (
	"attendance_session_service/internal/models"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type MainServiceClient struct {
	baseURL string
	client  *resty.Client
	token   string
}

func NewMainServiceClient(baseURL, token string, timeoutSeconds int) *MainServiceClient {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 10
	}

	return &MainServiceClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		client:  resty.New().SetTimeout(time.Duration(timeoutSeconds) * time.Second),
	}
}

func (c *MainServiceClient) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	resp, err := c.request(ctx).SetResult(&user).Get(c.baseURL + "/internal/users/" + strconv.Itoa(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get user from main service: %w", err)
	}
	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrRemoteNotFound
	}
	if resp.IsError() {
		return nil, fmt.Errorf("main service returned status %d for user %d", resp.StatusCode(), id)
	}
	return &user, nil
}

func (c *MainServiceClient) GetDepartmentByID(ctx context.Context, id int) (*models.Department, error) {
	var department models.Department
	resp, err := c.request(ctx).SetResult(&department).Get(c.baseURL + "/internal/departments/" + strconv.Itoa(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get department from main service: %w", err)
	}
	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrRemoteNotFound
	}
	if resp.IsError() {
		return nil, fmt.Errorf("main service returned status %d for department %d", resp.StatusCode(), id)
	}
	return &department, nil
}

func (c *MainServiceClient) GetDepartmentUsers(ctx context.Context, departmentID int) ([]models.User, error) {
	var users []models.User
	resp, err := c.request(ctx).SetResult(&users).Get(c.baseURL + "/internal/departments/" + strconv.Itoa(departmentID) + "/users")
	if err != nil {
		return nil, fmt.Errorf("failed to get department users from main service: %w", err)
	}
	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrRemoteNotFound
	}
	if resp.IsError() {
		return nil, fmt.Errorf("main service returned status %d for department users %d", resp.StatusCode(), departmentID)
	}
	return users, nil
}

func (c *MainServiceClient) request(ctx context.Context) *resty.Request {
	return c.client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetHeader("X-Service-Token", c.token)
}

var ErrRemoteNotFound = fmt.Errorf("remote resource not found")

package client

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type AuthClient struct {
	client  *resty.Client
	baseURL string
}

func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{
		client:  resty.New(),
		baseURL: baseURL,
	}
}

type ValidateResponse struct {
	Valid  bool `json:"valid"`
	UserID int  `json:"user_id"`
}

func (a *AuthClient) ValidateToken(token string) (*ValidateResponse, error) {
	var resp ValidateResponse

	r, err := a.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetResult(&resp).
		Get(a.baseURL + "/auth/validate")

	// 👇 ЛОГ 1: ошибка HTTP запроса
	if err != nil {
		fmt.Println("HTTP ERROR:", err)
		return nil, err
	}

	// 👇 ЛОГ 2: статус и body ответа
	fmt.Println("STATUS:", r.Status())
	fmt.Println("RAW RESPONSE:", r.String())

	// 👇 если статус 400/401/500
	if r.IsError() {
		return nil, fmt.Errorf("auth error: %s", r.String())
	}

	// 👇 ЛОГ 3: что распарсилось в struct
	fmt.Printf("PARSED RESPONSE: %+v\n", resp)

	return &resp, nil
}

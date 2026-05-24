package inputs

import "time"

type AuthUser struct {
	Id        int
	Email     string
	Password  string
	DeletedAt *time.Time
}

type RegisterInput struct {
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	DepartmentId *int   `json:"department_id"`
}
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UpdateUserProfileInput struct {
	FullName *string `json:"full_name"`
	Email    *string `json:"email"`
}

type UpdatePasswordInput struct {
	Password string `json:"password"`
}

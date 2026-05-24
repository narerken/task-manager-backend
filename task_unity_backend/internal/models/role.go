package models

type Role struct {
	Id           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	DepartmentId *int   `json:"department_id,omitempty" db:"department_id"`
}

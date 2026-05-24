package inputs

import "time"

type AddAttendanceInput struct {
	Date    time.Time `json:"date" db:"date"`
	UserId  int       `json:"user_id" db:"user_id"`
	Comment *string   `json:"comment,omitempty" db:"comment"`
	Creator int       `json:"-"`
}

type UpdateAttendanceInput struct {
	Status   *string `json:"status,omitempty" validate:"oneof=present absent excused"`
	Comment  *string `json:"comment"`
	MarkedBy *int    `json:"-"`
}

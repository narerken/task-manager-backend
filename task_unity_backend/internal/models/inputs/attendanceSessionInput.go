package inputs

type CreateAttendanceSessionInput struct {
	DepartmentId int    `json:"department_id" validate:"required"`
	Date         string `json:"date" validate:"required"`
}

type AttendanceEntryUpsertInput struct {
	StudentId int     `json:"student_id" validate:"required"`
	Status    string  `json:"status" validate:"required,oneof=present absent late excused"`
	Comment   *string `json:"comment,omitempty"`
}

type BulkUpsertAttendanceEntriesInput struct {
	Entries []AttendanceEntryUpsertInput `json:"entries"`
}

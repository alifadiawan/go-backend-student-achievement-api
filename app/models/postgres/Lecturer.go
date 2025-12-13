package postgres

type LecturerWithUser struct {
	ID         string `json:"id"`
	LecturerID string `json:"lecturer_id"`
	Department string `json:"department"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
}
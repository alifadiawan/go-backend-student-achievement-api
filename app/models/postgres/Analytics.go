package postgres

import "time"

type PeriodStat struct {
	Period time.Time `json:"period"`
	Total  int       `json:"total"`
}

type TopStudentStat struct {
	StudentID     string `json:"student_id"`
	TotalVerified int    `json:"total_verified"`
}

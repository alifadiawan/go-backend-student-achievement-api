package postgres


type StudentWithUser struct {
	ID           string
	StudentID    string
	ProgramStudy string
	AcademicYear string
	FullName     string
	Email        string
}

type StudentDetail struct {
	ID           string
	StudentID    string
	ProgramStudy string
	AcademicYear string
	FullName     string
	Email        string
	AdvisorID    *string
	AdvisorName  *string
}
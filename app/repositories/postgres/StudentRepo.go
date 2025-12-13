package postgres

import (
	models "backendUAS/app/models/postgres"
	"backendUAS/databases"
	"database/sql"
)

func GetAllStudentsRepo() ([]models.StudentWithUser, error) {
	rows, err := databases.DatabaseQuery.Query(`
		SELECT 
			s.id,
			s.student_id,
			s.program_study,
			s.academic_year,
			u.full_name,
			u.email
		FROM students s
		JOIN users u ON u.id = s.user_id
		ORDER BY u.full_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.StudentWithUser

	for rows.Next() {
		var s models.StudentWithUser
		if err := rows.Scan(
			&s.ID,
			&s.StudentID,
			&s.ProgramStudy,
			&s.AcademicYear,
			&s.FullName,
			&s.Email,
		); err != nil {
			return nil, err
		}
		results = append(results, s)
	}

	return results, nil
}

func GetStudentsByAdvisorRepo(lecturerID string) ([]models.StudentWithUser, error) {
	rows, err := databases.DatabaseQuery.Query(`
		SELECT 
			s.id,
			s.student_id,
			s.program_study,
			s.academic_year,
			u.full_name,
			u.email
		FROM students s
		JOIN users u ON u.id = s.user_id
		WHERE s.advisor_id = $1
		ORDER BY u.full_name
	`, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.StudentWithUser
	for rows.Next() {
		var s models.StudentWithUser
		if err := rows.Scan(
			&s.ID,
			&s.StudentID,
			&s.ProgramStudy,
			&s.AcademicYear,
			&s.FullName,
			&s.Email,
		); err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	return results, nil
}





func GetAchievementsByStudentIDRepo(studentID string) ([]models.Achievement, error) {
	rows, err := databases.DatabaseQuery.Query(`
		SELECT
			id,
			student_id,
			mongo_achievement_id,
			status,
			submitted_at,
			rejection_note,
			created_at,
			updated_at
		FROM achievement_references
		WHERE student_id = $1
		ORDER BY created_at DESC
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Achievement

	for rows.Next() {
		var a models.Achievement
		err := rows.Scan(
			&a.ID,
			&a.StudentID,
			&a.MongoId,
			&a.Status,
			&a.SubmittedAt,
			&a.RejectionNote,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, a)
	}

	return results, nil
}




func GetStudentByIDRepo(id string) (*models.StudentDetail, error) {
	row := databases.DatabaseQuery.QueryRow(`
		SELECT
			s.id,
			s.student_id,
			s.program_study,
			s.academic_year,
			u.full_name,
			u.email,
			l.id,
			ul.full_name
		FROM students s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN lecturers l ON l.id = s.advisor_id
		LEFT JOIN users ul ON ul.id = l.user_id
		WHERE s.id = $1
	`, id)

	var result models.StudentDetail

	err := row.Scan(
		&result.ID,
		&result.StudentID,
		&result.ProgramStudy,
		&result.AcademicYear,
		&result.FullName,
		&result.Email,
		&result.AdvisorID,
		&result.AdvisorName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &result, nil
}



func UpdateStudentAdvisorRepo(studentID, advisorID string) error {
	_, err := databases.DatabaseQuery.Exec(`
		UPDATE students
		SET advisor_id = $1
		WHERE id = $2
	`, advisorID, studentID)

	return err
}

func GetStudentAdvisorRepo(studentID string) (*models.LecturerWithUser, error) {
	row := databases.DatabaseQuery.QueryRow(`
		SELECT
			l.id,
			l.lecturer_id,
			l.department,
			u.full_name,
			u.email
		FROM students s
		JOIN lecturers l ON l.id = s.advisor_id
		JOIN users u ON u.id = l.user_id
		WHERE s.id = $1
	`, studentID)

	var advisor models.LecturerWithUser

	err := row.Scan(
		&advisor.ID,
		&advisor.LecturerID,
		&advisor.Department,
		&advisor.FullName,
		&advisor.Email,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &advisor, nil
}

func GetStudentAdviseesRepo(db *sql.DB) ([]models.LecturerWithUser, error) {
	rows, err := db.Query(`
		SELECT
			l.id,
			l.lecturer_id,
			l.department,
			u.full_name,
			u.email
		FROM lecturers l
		JOIN users u ON u.id = l.user_id
		JOIN roles r ON r.id = u.role_id
		WHERE r.name = 'dosen'
		ORDER BY u.full_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.LecturerWithUser

	for rows.Next() {
		var l models.LecturerWithUser
		if err := rows.Scan(
			&l.ID,
			&l.LecturerID,
			&l.Department,
			&l.FullName,
			&l.Email,
		); err != nil {
			return nil, err
		}
		results = append(results, l)
	}

	return results, nil
}

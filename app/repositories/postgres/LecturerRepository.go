package postgres

import (
	models "backendUAS/app/models/postgres"
	"backendUAS/databases"
)

func GetAllLecturerRepo() ([]models.LecturerWithUser, error) {
	rows, err := databases.DatabaseQuery.Query(`
		SELECT
			l.id,
			l.lecturer_id,
			l.department,
			u.full_name,
			u.email
		FROM lecturers l
		JOIN users u ON u.id = l.user_id
		ORDER BY u.full_name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lecturers []models.LecturerWithUser

	for rows.Next() {
		var l models.LecturerWithUser
		err := rows.Scan(
			&l.ID,
			&l.LecturerID,
			&l.Department,
			&l.FullName,
			&l.Email,
		)
		if err != nil {
			return nil, err
		}

		lecturers = append(lecturers, l)
	}

	return lecturers, nil
}

func GetLecturerByStudentIDRepo(studentID string) (*models.LecturerWithUser, error) {
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

	var lecturer models.LecturerWithUser

	err := row.Scan(
		&lecturer.ID,
		&lecturer.LecturerID,
		&lecturer.Department,
		&lecturer.FullName,
		&lecturer.Email,
	)
	if err != nil {
		return nil, err
	}

	return &lecturer, nil
}


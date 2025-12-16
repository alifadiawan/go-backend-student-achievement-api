package postgres

import (
	models "backendUAS/app/models/postgres"
	"backendUAS/databases"
)

func GetTotalAchievementByPeriodRepo() ([]models.PeriodStat, error) {
	rows, err := databases	.DatabaseQuery.Query(`
		SELECT
			DATE_TRUNC('month', verified_at) AS period,
			COUNT(*) AS total
		FROM achievement_references
		WHERE status = 'verified'
		GROUP BY period
		ORDER BY period
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.PeriodStat

	for rows.Next() {
		var stat models.PeriodStat
		err := rows.Scan(&stat.Period, &stat.Total)
		if err != nil {
			return nil, err
		}

		results = append(results, stat)
	}

	return results, nil
}


func GetTopStudentsRepo(limit int) ([]models.TopStudentStat, error) {
	rows, err := databases.DatabaseQuery.Query(`
		SELECT
			student_id,
			COUNT(*) AS total_verified
		FROM achievement_references
		WHERE status = 'verified'
		GROUP BY student_id
		ORDER BY total_verified DESC
		LIMIT $1
	`, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.TopStudentStat

	for rows.Next() {
		var stat models.TopStudentStat
		err := rows.Scan(&stat.StudentID, &stat.TotalVerified)
		if err != nil {
			return nil, err
		}

		results = append(results, stat)
	}

	return results, nil
}


func GetTotalVerifiedByStudentRepo(studentID string) (int, error) {
	var total int

	err := databases.DatabaseQuery.QueryRow(`
		SELECT COUNT(*)
		FROM achievement_references
		WHERE student_id = $1
		AND status = 'verified'
	`, studentID).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

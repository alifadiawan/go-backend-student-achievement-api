package postgres

import (
	models "backendUAS/app/models/postgres"
	"backendUAS/databases"
	"context"
	"database/sql"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var databaseName = os.Getenv("DB_NAME_MONGO")

func GetUserIDofAchievementRepo(achievement_references_id string) (string, error) {
	var studentID string

	query := `
		SELECT student_id
		FROM achievement_references
		WHERE id = $1
	`

	err := databases.DatabaseQuery.QueryRow(query, achievement_references_id).Scan(&studentID)
	if err != nil {
		return "", err
	}

	return studentID, nil
}

func GetAdvisorFromStudent(student_id string) (string, error) {

	var advisor_id string

	err := databases.DatabaseQuery.QueryRow(`
		SELECT advisor_id
		FROM students
		WHERE id = $1
	`, student_id).Scan(&advisor_id)

	if err != nil {
		return "", err
	}

	return advisor_id, err

}

func GetAllAchievementRepo() ([]models.Achievement, error) {
	rows, err := databases.DatabaseQuery.Query(`
		SELECT id, student_id, mongo_achievement_id, status, submitted_at,
			   rejection_note, created_at
		FROM achievement_references
		WHERE status != 'deleted'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Achievement

	for rows.Next() {
		var ac models.Achievement

		err := rows.Scan(
			&ac.ID,
			&ac.StudentID,
			&ac.MongoId,
			&ac.Status,
			&ac.SubmittedAt,
			&ac.RejectionNote,
			&ac.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, ac)
	}

	return results, nil
}

func GetAllAchievementByStudentIDRepo(studentID string) ([]models.Achievement, error) {
	rows, err := databases.DatabaseQuery.Query(`
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
			   rejection_note, created_at
		FROM achievement_references
		WHERE status != 'deleted'
		AND student_id = $1
	`, studentID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Achievement

	for rows.Next() {
		var ac models.Achievement
		err := rows.Scan(
			&ac.ID,
			&ac.StudentID,
			&ac.MongoId,
			&ac.Status,
			&ac.SubmittedAt,
			&ac.RejectionNote,
			&ac.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, ac)
	}

	return results, nil
}

func GetAchievementByIDRepo(StudentID string) ([]models.AchievementGabung, error) {

	rows, err := databases.DatabaseQuery.Query(`
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
			   rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE id = $1
	`, StudentID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.AchievementGabung

	collection := databases.MongoClient.Database(databaseName).Collection("achievements")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for rows.Next() {
		var ac models.Achievement

		err := rows.Scan(
			&ac.ID,
			&ac.StudentID,
			&ac.MongoId,
			&ac.Status,
			&ac.SubmittedAt,
			&ac.RejectionNote,
			&ac.CreatedAt,
			&ac.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Query MongoDB
		objID, _ := primitive.ObjectIDFromHex(ac.MongoId)

		var mongoData map[string]interface{}
		err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&mongoData)
		if err != nil {
			mongoData = make(map[string]interface{})
		}

		// Append ke hasil final
		results = append(results, models.AchievementGabung{
			Achievement: ac,
			Details:     mongoData,
		})
	}

	return results, nil
}

func InsertAchievementPostgres(studentID string, mongoID string, status string) error {
	res, err := databases.DatabaseQuery.Exec(`
		INSERT INTO achievement_references
			(student_id, mongo_achievement_id, status, submitted_at, created_at)
		VALUES ($1, $2, $3, now(), now())
	`, studentID, mongoID, status)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteAchievementRepo(achievement_references_id string) (bool, error) {

	result, err := databases.DatabaseQuery.Exec(`
		UPDATE achievement_references
			SET status = 'deleted'
		WHERE id = $1
	`, achievement_references_id)
	if err != nil {
		return false, err
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsEffected == 0 {
		return false, nil
	}

	return true, err

}

func SubmitAchievementRepository(AchievementID string, studentID string) (bool, error) {

	query, err := databases.DatabaseQuery.Exec(`
		UPDATE achievement_references 
		SET status = 'submitted'
		WHERE id = $1 AND student_id = $2
	`, AchievementID, studentID)

	if err != nil {
		return false, err
	}

	rowsEffected, err := query.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsEffected == 0 {
		return false, err
	}

	return true, err

}

func ApproveAchievmentRepository(AchievmentID string) (bool, error) {

	query, err := databases.DatabaseQuery.Exec(`
		UPDATE achievment_references
			SET status = 'approved'
		WHERE id = $1
	`)

	if err != nil {
		return false, err
	}

	rowsEffected, err := query.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsEffected == 0 {
		return false, err
	}

	return true, err

}

func VerifyAchievementRepo(achievment_references_id string) (bool, error) {

	query, err := databases.DatabaseQuery.Exec(`
		UPDATE achievement_references
			SET status = 'verified'
		WHERE id = $1
	`, achievment_references_id)

	if err != nil {
		return false, err
	}

	result, err := query.RowsAffected()
	if err != nil {
		return false, err
	}

	if result == 0 {
		return false, nil
	}

	return true, err

}

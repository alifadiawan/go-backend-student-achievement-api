package postgres

import (
	models "backendUAS/app/models/postgres"
	"backendUAS/databases"
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var databaseName = os.Getenv("DB_NAME_MONGO")

func GetAllAchievementRepo() ([]models.AchievementGabung, error) {

	rows, err := databases.DatabaseQuery.Query(`
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
			   rejection_note, created_at, updated_at
		FROM achievement_references
	`)
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

		objID, objErr := primitive.ObjectIDFromHex(ac.MongoId)
		if objErr != nil {
			continue 
		}

		var mongoData map[string]interface{}
		err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&mongoData)
		if err != nil {
			fmt.Println("Mongo Find Error:", err)
			return nil, err
		}

		results = append(results, models.AchievementGabung{
			Achievement: ac,
			Details:     mongoData,
		})
	}

	return results, nil
}

func GetAchievementByIDRepo(studentID string) ([]models.AchievementGabung, error) {

	rows, err := databases.DatabaseQuery.Query(`
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
			   rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1
	`, studentID)

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

package mongo

import (
	models "backendUAS/app/models/mongo"
	"backendUAS/databases"
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var databaseName = os.Getenv("DB_NAME_MONGO")

func GetAllAchievementByIDRepo(userid string) ([]models.Achievement, error) {

	collection := databases.MongoClient.Database(databaseName).Collection("student_achievement")

	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	filter := bson.M{}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var AllAchievement []models.Achievement
	if err := cursor.All(ctx, &AllAchievement); err != nil {
		return nil, err
	}

	return AllAchievement, err

}

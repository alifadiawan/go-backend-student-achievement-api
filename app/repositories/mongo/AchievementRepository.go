package mongo

import (
	models "backendUAS/app/models/mongo"
	"backendUAS/databases"
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var databaseName = os.Getenv("DB_NAME_MONGO")

func GetAllAchievementByIDRepo(userid string) ([]models.Achievement, error) {

	collection := databases.MongoClient.Database(databaseName).Collection("student_achievement")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

func AddAchievementRepositoryMongo(achievement models.Achievement) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	result, err := databases.MongoClient.Database("student_achievement").
		Collection("achievements").InsertOne(ctx, achievement)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

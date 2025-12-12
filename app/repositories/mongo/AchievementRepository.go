package mongo

import (
	models "backendUAS/app/models/mongo"
	"backendUAS/databases"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func UploadAchievementRepo(achievementRequest models.AchievementAttachementRequest) (bool, error) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if achievementRequest.Attachment == nil {
		return false, fmt.Errorf("no file provided")
	}

	if achievementRequest.AchievementReferencesID == "" {
		return false, fmt.Errorf("achievement_references_id is required")
	}

	basePath := "achievement_attachment"
	targetFolder := filepath.Join(basePath, achievementRequest.AchievementReferencesID)

	err := os.MkdirAll(targetFolder, os.ModePerm)
	if err != nil {
		return false, fmt.Errorf("failed to create folder: %v", err)
	}

	// Open uploaded file
	file, err := achievementRequest.Attachment.Open()
	if err != nil {
		return false, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create unique filename
	ext := filepath.Ext(achievementRequest.Attachment.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Destination path
	destPath := filepath.Join(targetFolder, fileName)

	// Create destination file
	dst, err := os.Create(destPath)
	if err != nil {
		return false, fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return false, fmt.Errorf("failed to save file: %v", err)
	}

	return true, nil
}


package mongo

import (
	"backendUAS/databases"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAchievementDetailByIdRepo(mongoID string) (map[string]interface{}, error) {

	objID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return nil, err
	}

	collection := databases.MongoClient.
		Database("student_achievement").
		Collection("achievements")

	var result bson.M

	err = collection.FindOne(
		context.Background(),
		bson.M{"_id": objID},
	).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

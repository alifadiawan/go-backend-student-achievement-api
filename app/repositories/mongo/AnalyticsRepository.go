package mongo

import (
	"backendUAS/databases"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTotalAchievementByTypeRepo() (map[string]int, error) {
	collection := databases.MongoClient.
		Database("student_achievement").
		Collection("achievements")

	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":   "$achievementType",
			"total": bson.M{"$sum": 1},
		}}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	results := make(map[string]int)

	for cursor.Next(context.TODO()) {
		var row struct {
			ID    string `bson:"_id"`
			Total int    `bson:"total"`
		}

		err := cursor.Decode(&row)
		if err != nil {
			return nil, err
		}

		results[row.ID] = row.Total
	}

	return results, nil
}


func GetCompetitionLevelDistributionRepo() (map[string]int, error) {
	collection := databases.MongoClient.
		Database("nama_database_kamu").
		Collection("achievements")

	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":   "$details.level",
			"total": bson.M{"$sum": 1},
		}}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	results := make(map[string]int)

	for cursor.Next(context.TODO()) {
		var row struct {
			ID    string `bson:"_id"`
			Total int    `bson:"total"`
		}

		if err := cursor.Decode(&row); err != nil {
			return nil, err
		}

		results[row.ID] = row.Total
	}

	return results, nil
}



func GetStudentAchievementByTypeRepo(studentID string) (map[string]int, error) {

	collection := databases.MongoClient.
		Database("nama_database_kamu").
		Collection("achievements")

	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"studentId": studentID,
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$achievementType",
			"total": bson.M{"$sum": 1},
		}}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	results := make(map[string]int)

	for cursor.Next(context.TODO()) {
		var row struct {
			ID    string `bson:"_id"`
			Total int    `bson:"total"`
		}

		err := cursor.Decode(&row)
		if err != nil {
			return nil, err
		}

		results[row.ID] = row.Total
	}

	return results, nil
}

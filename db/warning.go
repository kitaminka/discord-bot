package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	WarningCollectionName = "warnings"
	WarningDuration       = 36 * time.Hour
)

type Warning struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Time        time.Time          `bson:"time,omitempty"`
	Reason      string             `bson:"reason,omitempty"`
	UserID      string             `bson:"userId,omitempty"`
	ModeratorID string             `bson:"moderatorId,omitempty"`
}

func CreateWarning(warn Warning) error {
	_, err := MongoDatabase.Collection(WarningCollectionName).InsertOne(context.Background(), warn)
	return err
}
func RemoveWarning(warnID primitive.ObjectID) (Warning, error) {
	var warning Warning
	err := MongoDatabase.Collection(WarningCollectionName).FindOneAndDelete(context.Background(), bson.D{{"_id", warnID}}).Decode(&warning)
	return warning, err
}
func RemoveUserWarnings(userID string) error {
	_, err := MongoDatabase.Collection(WarningCollectionName).DeleteMany(context.Background(), bson.D{{"userId", userID}})
	return err
}
func GetWarning(warnID uint64) (Warning, error) {
	var warn Warning
	err := MongoDatabase.Collection(WarningCollectionName).FindOne(context.Background(), bson.D{{"id", warnID}}).Decode(&warn)
	return warn, err
}
func GetUserWarnings(userID string) ([]Warning, error) {
	err := RemoveExpiredUserWarnings(userID)
	if err != nil {
		return nil, err
	}
	var warnings []Warning
	cursor, err := MongoDatabase.Collection(WarningCollectionName).Find(context.Background(), bson.D{{"userId", userID}})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &warnings)
	return warnings, err
}
func RemoveExpiredUserWarnings(userID string) error {
	_, err := MongoDatabase.Collection(WarningCollectionName).DeleteMany(context.Background(), bson.D{{"userId", userID}, {"time", bson.D{{"$lt", time.Now().Add(-WarningDuration)}}}})
	return err
}

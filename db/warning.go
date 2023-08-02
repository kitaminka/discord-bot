package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	WarningCollectionName = "warnings"
	WarnDuration          = 36 * time.Hour
)

type Warning struct {
	ID          uint64    `bson:"id,omitempty"`
	Time        time.Time `bson:"time,omitempty"`
	UserID      string    `bson:"userId,omitempty"`
	ModeratorID string    `bson:"moderatorId,omitempty"`
}

func AddUserWarning(warn Warning) error {
	_, err := MongoDatabase.Collection(WarningCollectionName).InsertOne(context.Background(), warn)
	return err
}
func RemoveWarning(warnID uint64) error {
	_, err := MongoDatabase.Collection(WarningCollectionName).DeleteOne(context.Background(), bson.D{{"id", warnID}})
	return err
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
	var warnings []Warning
	cursor, err := MongoDatabase.Collection(WarningCollectionName).Find(context.Background(), bson.D{{"userId", userID}})
	if err != nil {
		return warnings, err
	}
	err = cursor.All(context.Background(), &warnings)
	return warnings, err
}

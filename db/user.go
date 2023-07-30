package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const UserCollectionName = "users"

const (
	ReputationDelay   = 40 * time.Minute
	DefaultReputation = 0
	WarnDuration      = 36 * time.Hour
)

var (
	MaxReputation float64 = 1000000
	MinReputation float64 = -1000000
)

type Warn struct {
	Time        time.Time `bson:"time,omitempty"`
	ModeratorID string    `bson:"moderatorId,omitempty"`
}
type User struct {
	ID               string    `bson:"id,omitempty"`
	Reputation       int       `bson:"reputation,omitempty"`
	ReputationDelay  time.Time `bson:"reputationDelayEnd,omitempty"`
	ReportsSentCount int       `bson:"reportsSentCount,omitempty"`
	Warns            []Warn    `bson:"warns,omitempty"`
}

func GetUser(userID string) (User, error) {
	var user User
	err := MongoDatabase.Collection(UserCollectionName).FindOne(context.Background(), bson.D{{"id", userID}}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return User{
			ID:               userID,
			ReputationDelay:  time.Time{},
			Reputation:       DefaultReputation,
			ReportsSentCount: 0,
		}, nil
	}
	return user, err
}

func RemoveUser(userID string) error {
	_, err := MongoDatabase.Collection(UserCollectionName).DeleteOne(context.Background(), bson.D{{"id", userID}})
	return err
}

func SetUserReputation(userID string, reputation int) error {
	_, err := MongoDatabase.Collection(UserCollectionName).UpdateOne(context.Background(), bson.D{{"id", userID}}, bson.D{{"$set", bson.D{{"reputation", reputation}}}}, options.Update().SetUpsert(true))
	return err
}
func ChangeUserReputation(userID string, change int) error {
	_, err := MongoDatabase.Collection(UserCollectionName).UpdateOne(context.Background(), bson.D{{"id", userID}}, bson.D{{"$inc", bson.D{{"reputation", change}}}}, options.Update().SetUpsert(true))
	return err
}
func UpdateUserReputationDelay(userID string) error {
	delayEnd := time.Now().Add(ReputationDelay)
	_, err := MongoDatabase.Collection(UserCollectionName).UpdateOne(context.Background(), bson.D{{"id", userID}}, bson.D{{"$set", bson.D{{"reputationDelayEnd", delayEnd}}}}, options.Update().SetUpsert(true))
	return err
}
func ResetUserReputationDelay(userID string) error {
	_, err := MongoDatabase.Collection(UserCollectionName).UpdateOne(context.Background(), bson.D{{"id", userID}}, bson.D{{"$unset", bson.D{{"reputationDelayEnd", nil}}}}, options.Update().SetUpsert(true))
	return err
}
func GetUserReputationTop() (*[]User, error) {
	users := &[]User{}
	aggregate, err := MongoDatabase.Collection(UserCollectionName).Aggregate(context.Background(), mongo.Pipeline{
		{{"$match", bson.D{{"reputation", bson.D{{"$exists", true}}}}}},
		{{"$sort", bson.D{{"reputation", -1}}}},
		{{"$limit", 10}},
	})
	if err != nil {
		return users, err
	}
	err = aggregate.All(nil, users)
	return users, err
}

func IncrementUserReportsSent(userID string) error {
	_, err := MongoDatabase.Collection(UserCollectionName).UpdateOne(context.Background(), bson.D{{"id", userID}}, bson.D{{"$inc", bson.D{{"reportsSentCount", 1}}}}, options.Update().SetUpsert(true))
	return err
}

func AddUserWarn(userID string, warn Warn) error {
	_, err := MongoDatabase.Collection(UserCollectionName).UpdateOne(context.Background(), bson.D{{"id", userID}}, bson.D{{"$push", bson.D{{"warns", warn}}}}, options.Update().SetUpsert(true))
	return err
}
func RemoveUserWarn(userID string, warn Warn) error {
	_, err := MongoDatabase.Collection(UserCollectionName).UpdateOne(context.Background(), bson.D{{"id", userID}}, bson.D{{"$pull", bson.D{{"warns", warn}}}}, options.Update().SetUpsert(true))
	return err
}

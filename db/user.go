package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const UserCollectionName = "users"

const ReputationDelay = 40 * time.Minute

var MaxReputation float64 = 1000000
var MinReputation float64 = -1000000

type User struct {
	ID               string    `bson:"id,omitempty"`
	Reputation       int       `bson:"reputation,omitempty"`
	ReputationDelay  time.Time `bson:"reputationDelayEnd,omitempty"`
	ReportsSentCount int       `bson:"reportsSentCount,omitempty"`
}

func GetUser(userID string) (User, error) {
	var user User
	err := MongoDatabase.Collection(UserCollectionName).FindOne(nil, bson.D{{"id", userID}}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return User{
			ID:               userID,
			ReputationDelay:  time.Time{},
			Reputation:       0,
			ReportsSentCount: 0,
		}, nil
	}
	return user, err
}

func SetUserReputation(userID string, reputation int) error {
	_, err := MongoDatabase.Collection(UserCollectionName).UpdateOne(nil, bson.D{{"id", userID}}, bson.D{{"$set", bson.D{{"reputation", reputation}}}}, options.Update().SetUpsert(true))
	return err
}
func ChangeUserReputation(userID string, change int) error {
	_, err := MongoDatabase.Collection(UserCollectionName).UpdateOne(nil, bson.D{{"id", userID}}, bson.D{{"$inc", bson.D{{"reputation", change}}}}, options.Update().SetUpsert(true))
	return err
}
func UpdateUserReputationDelay(userID string) error {
	delayEnd := time.Now().Add(ReputationDelay)
	_, err := MongoDatabase.Collection(UserCollectionName).UpdateOne(nil, bson.D{{"id", userID}}, bson.D{{"$set", bson.D{{"reputationDelayEnd", delayEnd}}}}, options.Update().SetUpsert(true))
	return err
}
func ResetUserReputationDelay(userID string) error {
	_, err := MongoDatabase.Collection(UserCollectionName).UpdateOne(nil, bson.D{{"id", userID}}, bson.D{{"$set", bson.D{{"reputationDelayEnd", time.Time{}}}}}, options.Update().SetUpsert(true))
	return err
}
func GetUserReputationTop() (*[]User, error) {
	users := &[]User{}
	aggregate, err := MongoDatabase.Collection(UserCollectionName).Aggregate(nil, mongo.Pipeline{
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
	_, err := MongoDatabase.Collection(UserCollectionName).UpdateOne(nil, bson.D{{"id", userID}}, bson.D{{"$inc", bson.D{{"reportsSentCount", 1}}}}, options.Update().SetUpsert(true))
	return err
}

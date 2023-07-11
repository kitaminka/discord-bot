package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const UserCollectionName = "users"

const ReputationDelay = 12 * time.Hour

type User struct {
	ID              string    `bson:"id,omitempty"`
	Reputation      int       `bson:"reputation,omitempty"`
	ReputationDelay time.Time `bson:"reputationDelayEnd,omitempty"`
}

func GetUser(userID string) (User, error) {
	var member User
	err := MongoDatabase.Collection(UserCollectionName).FindOne(nil, bson.D{{"id", userID}}).Decode(&member)
	if err == mongo.ErrNoDocuments {
		return User{
			ID:              userID,
			ReputationDelay: time.Time{},
			Reputation:      0,
		}, nil
	}
	return member, err
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

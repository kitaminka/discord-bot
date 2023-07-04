package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MembersCollectionName = "members"

type Member struct {
	ID         string `bson:"id,omitempty"`
	Reputation int    `bson:"reputation,omitempty"`
}

func GetMember(memberID string) (Member, error) {
	var member Member
	err := MongoDatabase.Collection(MembersCollectionName).FindOne(nil, bson.D{{"id", memberID}}).Decode(&member)
	if err == mongo.ErrNoDocuments {
		return Member{
			ID:         memberID,
			Reputation: 0,
		}, nil
	}
	return member, err
}

func ChangeMemberReputation(memberID string, change int) error {
	_, err := MongoDatabase.Collection(MembersCollectionName).UpdateOne(nil, bson.D{{"id", memberID}}, bson.D{{"$inc", bson.D{{"reputation", change}}}}, options.Update().SetUpsert(true))
	return err
}

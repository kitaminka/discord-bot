package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const GuildsCollectionName = "guild"

// Only main guild in this collection

type Guild struct {
	ID                     string `bson:"id,omitempty"`
	ReportChannelID        string `bson:"reportChannelId,omitempty"`
	ResoledReportChannelID string `bson:"resoledReportChannelId,omitempty"`
	ReputationLogChannelID string `bson:"reputationLogChannelId,omitempty"`
}

func GetGuild() (Guild, error) {
	var server Guild
	err := MongoDatabase.Collection(GuildsCollectionName).FindOne(context.Background(), bson.D{}).Decode(&server)
	return server, err
}
func UpdateGuild(server Guild) error {
	_, err := MongoDatabase.Collection(GuildsCollectionName).UpdateOne(context.Background(), bson.D{}, bson.D{{"$set", server}}, options.Update().SetUpsert(true))
	return err
}

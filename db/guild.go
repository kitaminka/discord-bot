package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const GuildCollectionName = "guild"

// Only main guild is stored in this collection

type Guild struct {
	ID                     string `bson:"id,omitempty"`
	SupremeModeratorRoleID string `bson:"supremeModeratorRoleID,omitempty"`
	ReportChannelID        string `bson:"reportChannelId,omitempty"`
	ResoledReportChannelID string `bson:"resoledReportChannelId,omitempty"`
	ReputationLogChannelID string `bson:"reputationLogChannelId,omitempty"`
	ModerationLogChannelID string `bson:"moderationLogChannelId,omitempty"`
}

func GetGuild() (Guild, error) {
	var server Guild
	err := MongoDatabase.Collection(GuildCollectionName).FindOne(context.Background(), bson.D{}).Decode(&server)
	return server, err
}
func UpdateGuild(server Guild) error {
	_, err := MongoDatabase.Collection(GuildCollectionName).UpdateOne(context.Background(), bson.D{}, bson.D{{"$set", server}}, options.Update().SetUpsert(true))
	return err
}

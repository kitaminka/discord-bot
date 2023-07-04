package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Guild struct {
	ID                     string `bson:"id,omitempty"`
	ReportChannelID        string `bson:"reportChannelId,omitempty"`
	ResoledReportChannelID string `bson:"resoledReportChannelId,omitempty"`
}

func UpdateGuild(server Guild) error {
	_, err := MongoDatabase.Collection("servers").UpdateOne(nil, bson.D{{"id", server.ID}}, bson.D{{"$set", server}}, options.Update().SetUpsert(true))
	return err
}

func GetGuild(discordId string) (Guild, error) {
	var server Guild
	err := MongoDatabase.Collection("servers").FindOne(nil, bson.D{{"id", discordId}}).Decode(&server)
	return server, err
}

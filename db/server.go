package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	ID                     string `bson:"id"`
	ReportChannelID        string `bson:"reportChannelId"`
	ResoledReportChannelId string `bson:"resoledReportChannelId"`
}

func UpdateServer(server Server) error {
	_, err := MongoDatabase.Collection("servers").UpdateOne(nil, bson.D{{"id", server.ID}}, bson.D{{"$set", server}}, options.Update().SetUpsert(true))
	return err
}

func GetServer(discordId string) (Server, error) {
	var server Server
	err := MongoDatabase.Collection("servers").FindOne(nil, bson.D{{"id", discordId}}).Decode(&server)
	return server, err
}

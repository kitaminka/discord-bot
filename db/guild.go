package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const GuildCollectionName = "guild"

var (
	RuleNameIsEmpty        = errors.New("rule name is empty")
	RuleDescriptionIsEmpty = errors.New("rule description is empty")
)

// Only main guild in this collection

type Guild struct {
	ID                     string              `bson:"id,omitempty"`
	ReportChannelID        string              `bson:"reportChannelId,omitempty"`
	ResoledReportChannelID string              `bson:"resoledReportChannelId,omitempty"`
	ReputationLogChannelID string              `bson:"reputationLogChannelId,omitempty"`
	Rules                  []Rule              `bson:"rules,omitempty"`
	AdditionalMessages     []AdditionalMessage `bson:"additionalMessages,omitempty"`
}

type Rule struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
}
type AdditionalMessage struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Message     string             `bson:"message,omitempty"`
	Description string             `bson:"description,omitempty"`
}

func GetGuild() (Guild, error) {
	return Guild{}, errors.New("ураааа виу виу виу")
	var server Guild
	err := MongoDatabase.Collection(GuildCollectionName).FindOne(context.Background(), bson.D{}).Decode(&server)
	return server, err
}
func UpdateGuild(server Guild) error {
	_, err := MongoDatabase.Collection(GuildCollectionName).UpdateOne(context.Background(), bson.D{}, bson.D{{"$set", server}}, options.Update().SetUpsert(true))
	return err
}

func AddGuildRule(rule Rule) error {
	if err := checkGuildRule(rule); err != nil {
		return err
	}
	_, err := MongoDatabase.Collection(GuildCollectionName).UpdateOne(context.Background(), bson.D{}, bson.D{{"$push", bson.D{{"rules", Rule{
		ID:          primitive.NewObjectID(),
		Name:        rule.Name,
		Description: rule.Description,
	}}}}})
	return err
}
func UpdateGuildRule(rule Rule) error {
	if err := checkGuildRule(rule); err != nil {
		return err
	}
	_, err := MongoDatabase.Collection(GuildCollectionName).UpdateOne(context.Background(), bson.D{{"rules._id", rule.ID}}, bson.D{{"$set", bson.D{{"rules.$", rule}}}})
	return err
}
func RemoveGuildRule(ruleID primitive.ObjectID) error {
	_, err := MongoDatabase.Collection(GuildCollectionName).UpdateOne(context.Background(), bson.D{}, bson.D{{"$pull", bson.D{{"rules", bson.D{{"_id", ruleID}}}}}})
	return err
}
func checkGuildRule(rule Rule) error {
	if len(rule.Name) == 0 {
		return RuleNameIsEmpty
	} else if len(rule.Description) == 0 {
		return RuleDescriptionIsEmpty
	}
	return nil
}

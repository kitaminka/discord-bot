package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ReasonCollectionName = "reasons"

var (
	ReasonNameIsEmpty        = errors.New("rule name is empty")
	ReasonDescriptionIsEmpty = errors.New("rule description is empty")
)

type Reason struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
}

func CreateReason(reason Reason) error {
	if err := checkReason(reason); err != nil {
		return err
	}
	_, err := MongoDatabase.Collection(ReasonCollectionName).InsertOne(context.Background(), reason)
	return err
}
func UpdateReason(reason Reason) error {
	if err := checkReason(reason); err != nil {
		return err
	}
	_, err := MongoDatabase.Collection(ReasonCollectionName).UpdateOne(context.Background(), bson.D{{"_id", reason.ID}}, bson.D{{"$set", reason}})
	return err
}
func DeleteReason(reasonID primitive.ObjectID) error {
	_, err := MongoDatabase.Collection(ReasonCollectionName).DeleteOne(context.Background(), bson.D{{"_id", reasonID}})
	return err
}
func GetReasons() ([]Reason, error) {
	var reasons []Reason
	cursor, err := MongoDatabase.Collection(ReasonCollectionName).Find(context.Background(), bson.D{})
	if err != nil {
		return reasons, err
	}
	err = cursor.All(context.Background(), &reasons)
	return reasons, err
}
func checkReason(rule Reason) error {
	if len(rule.Name) == 0 {
		return ReasonNameIsEmpty
	} else if len(rule.Description) == 0 {
		return ReasonDescriptionIsEmpty
	}
	return nil
}

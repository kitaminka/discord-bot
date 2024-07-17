package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func CheckDbSpeed() {
	start := time.Now()
	var user User
	_ = MongoDatabase.Collection(UserCollectionName).FindOne(context.Background(), bson.D{{"id", "890320305082478652"}}).Decode(&user)
	timeElapsed := time.Since(start)
	fmt.Printf("took %s\n", timeElapsed)
	fmt.Println(user)
}

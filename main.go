package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kitaminka/discord-bot/bot"
)

var (
	Token             string
	MongoUri          string
	MongoDatabaseName string
)

func init() {
	godotenv.Load()

	Token = os.Getenv("DISCORD_TOKEN")
	MongoUri = os.Getenv("MONGO_URI")
	MongoDatabaseName = os.Getenv("MONGO_DATABASE_NAME")
}

func main() {
	bot.StartBot(Token, MongoUri, MongoDatabaseName)
}

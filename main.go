package main

import (
	"github.com/joho/godotenv"
	"github.com/kitaminka/discord-bot/bot"
	"log"
	"os"
)

var (
	Token    string
	MongoUri string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("Error loading .env file: %v", err)
	}

	Token = os.Getenv("DISCORD_TOKEN")
	MongoUri = os.Getenv("MONGO_URI")
}

func main() {
	bot.StartBot(Token, MongoUri)
}

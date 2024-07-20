package bot

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/interactions"
	"github.com/kitaminka/discord-bot/setup"
)

const (
	Intents = 1535
)

func StartBot(token, mongoUri, mongoDatabaseName string) {
	db.ConnectMongo(mongoUri, mongoDatabaseName)

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Panicf("Error creating Discord session: %v", err)
	}

	AddHandlers(session)
	go interactions.IntervalDeleteExpiredWarnings()
	session.Identify.Intents = Intents

	err = session.Open()
	if err != nil {
		log.Panicf("Error opening Discord session: %v", err)
	}

	setup.CreateSetupCommand(session)

	log.Print("Bot is now running. Press CTRL-C to exit.")

	defer session.Close()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}

package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/interactions"
	"log"
	"os"
	"os/signal"
)

const (
	Intents = 1535
)

func StartBot(token, mongoUri string) {
	db.ConnectMongo(mongoUri)
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Panicf("Error creating Discord session: %v", err)
	}

	AddHandlers(session)
	session.Identify.Intents = Intents

	err = session.Open()
	if err != nil {
		log.Panicf("Error opening Discord session: %v", err)
	}

	interactions.RemoveApplicationCommands(session)
	interactions.CreateApplicationCommands(session)

	log.Println("Bot is now running. Press CTRL-C to exit.")

	defer session.Close()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}

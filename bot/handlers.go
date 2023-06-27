package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/interactions"
)

var Handlers = []interface{}{
	func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
		switch interactionCreate.Type {
		case discordgo.InteractionApplicationCommand:
			interactions.Commands[interactionCreate.ApplicationCommandData().Name].Handler(session, interactionCreate)
		case discordgo.InteractionMessageComponent:
			interactions.ComponentHandlers[interactionCreate.MessageComponentData().CustomID](session, interactionCreate)
		}
	},
}

func AddHandlers(session *discordgo.Session) {
	for _, value := range Handlers {
		session.AddHandler(value)
	}
}

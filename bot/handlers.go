package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/cmd"
)

var Handlers = []interface{}{
	func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
		switch interactionCreate.Type {
		case discordgo.InteractionApplicationCommand:
			cmd.Commands[interactionCreate.ApplicationCommandData().Name].Handler(session, interactionCreate)
		}
	},
}

func AddHandlers(session *discordgo.Session) {
	for _, value := range Handlers {
		session.AddHandler(value)
	}
}

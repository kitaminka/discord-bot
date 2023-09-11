package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/interactions"
)

var Handlers = []interface{}{
	func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
		switch interactionCreate.Type {
		case discordgo.InteractionApplicationCommand:
			interactions.CommandHandlers[interactionCreate.ApplicationCommandData().Name](session, interactionCreate)
		case discordgo.InteractionMessageComponent:
			interactions.ComponentHandlers[interactionCreate.MessageComponentData().CustomID](session, interactionCreate)
		case discordgo.InteractionApplicationCommandAutocomplete:
			interactions.AutocompleteHandlers[interactionCreate.ApplicationCommandData().Name](session, interactionCreate)
		}
	},
	func(session *discordgo.Session, guildMemberRemove *discordgo.GuildMemberRemove) {
		err := db.RemoveUser(guildMemberRemove.User.ID)
		if err != nil {
			fmt.Printf("Error removing user: %v", err)
		}
	},
}

func AddHandlers(session *discordgo.Session) {
	for _, value := range Handlers {
		session.AddHandler(value)
	}
}

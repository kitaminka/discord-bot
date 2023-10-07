package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/interactions"
	"strings"
)

var Handlers = []interface{}{
	func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
		switch interactionCreate.Type {
		case discordgo.InteractionApplicationCommand:
			handler, exists := interactions.CommandHandlers[interactionCreate.ApplicationCommandData().Name]
			if !exists {
				interactions.InteractionRespondError(session, interactionCreate.Interaction, "Команда не найдена. Свяжитесь с администрацией.")
				return
			}
			handler(session, interactionCreate)
		case discordgo.InteractionMessageComponent:
			customID := strings.Split(interactionCreate.MessageComponentData().CustomID, ":")[0]
			handler, exists := interactions.ComponentHandlers[customID]
			if !exists {
				interactions.InteractionRespondError(session, interactionCreate.Interaction, "Команда не найдена. Свяжитесь с администрацией.")
				return
			}
			handler(session, interactionCreate)
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

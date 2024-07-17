package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/interactions"
)

var Handlers = []interface{}{
	func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
		switch interactionCreate.Type {
		case discordgo.InteractionApplicationCommand:
			if interactionCreate.ApplicationCommandData().Name == "setup" {
				interactions.SetupCommandHandler(session, interactionCreate)
				return
			}
			handler, exists := interactions.CommandHandlers[interactionCreate.ApplicationCommandData().Name]
			if !exists {
				interactions.InteractionRespondError(session, interactionCreate.Interaction, "Команда не найдена. Свяжитесь с администрацией.")
				return
			}
			handler(session, interactionCreate)
		case discordgo.InteractionMessageComponent:
			// Component name format: <name>:<someID>:<someID>:<someID>
			customID := strings.Split(interactionCreate.MessageComponentData().CustomID, ":")[0]
			handler, exists := interactions.ComponentHandlers[customID]
			if !exists {
				interactions.InteractionRespondError(session, interactionCreate.Interaction, "Команда не найдена. Свяжитесь с администрацией.")
				return
			}
			handler(session, interactionCreate)
		}
	},
	// func(session *discordgo.Session, guildMemberRemove *discordgo.GuildMemberRemove) {
	// 	err := db.RemoveUser(guildMemberRemove.User.ID)
	// 	if err != nil {
	// 		fmt.Printf("Error removing user: %v", err)
	// 	}
	// },
}

func AddHandlers(session *discordgo.Session) {
	for _, value := range Handlers {
		session.AddHandler(value)
	}
}

package interactions

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type Command struct {
	ApplicationCommand *discordgo.ApplicationCommand
	Handler            func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
}

var Commands = map[string]Command{
	"Report": {
		ApplicationCommand: &discordgo.ApplicationCommand{
			Type: discordgo.MessageApplicationCommand,
			Name: "Report",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian: "Отправить репорт",
			},
		},
		Handler: reportMessageCommandHandler,
	},
	"report": {
		ApplicationCommand: &discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "report",
			Description: "Отправить репорт",
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Report",
				},
			})
		},
	},
}

func CreateApplicationCommands(session *discordgo.Session) {
	for index, value := range Commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "1096521857081036831", value.ApplicationCommand)
		if err != nil {
			log.Panicf("Error creating '%v' command: %v", value.ApplicationCommand.Name, err)
		}
		log.Printf("Successfully created '%v' command", cmd.Name)

		if command, exists := Commands[index]; exists {
			command.ApplicationCommand = cmd
			Commands[index] = command
		}
	}
}
func RemoveApplicationCommands(session *discordgo.Session) {
	for _, value := range Commands {
		err := session.ApplicationCommandDelete(session.State.User.ID, "1096521857081036831", value.ApplicationCommand.ID)
		if err != nil {
			log.Panicf("Error deleting '%v' command: %v", value.ApplicationCommand.Name, err)
		}
		log.Printf("Successfully deleted '%v' command", value.ApplicationCommand.Name)
	}
}

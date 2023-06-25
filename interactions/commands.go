package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/cfg"
	"log"
)

type Command struct {
	ApplicationCommand *discordgo.ApplicationCommand
	Handler            func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
}

var Commands = map[string]Command{
	"Отправить репорт": {
		ApplicationCommand: &discordgo.ApplicationCommand{
			Type: discordgo.MessageApplicationCommand,
			Name: "Отправить репорт",
		},
		Handler: reportMessageCommandHandler,
	},
}

func CreateApplicationCommands(session *discordgo.Session) {
	for index, value := range Commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, cfg.Config.ServerID, value.ApplicationCommand)
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
		err := session.ApplicationCommandDelete(session.State.User.ID, cfg.Config.ServerID, value.ApplicationCommand.ID)
		if err != nil {
			log.Panicf("Error deleting '%v' command: %v", value.ApplicationCommand.Name, err)
		}
		log.Printf("Successfully deleted '%v' command", value.ApplicationCommand.Name)
	}
}

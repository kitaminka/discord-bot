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
	commands, err := session.ApplicationCommands(session.State.User.ID, cfg.Config.ServerID)
	if err != nil {
		log.Printf("Error getting application commands: %v", err)
		return
	}
	for _, command := range commands {
		err = session.ApplicationCommandDelete(session.State.User.ID, cfg.Config.ServerID, command.ID)
		if err != nil {
			log.Panicf("Error deleting '%v' command: %v", command.Name, err)
		}
		log.Printf("Successfully deleted '%v' command", command.Name)
	}
}

package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"log"
)

type Command struct {
	ApplicationCommand *discordgo.ApplicationCommand
	Handler            func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
}

var (
	Commands = map[string]Command{
		"Отправить репорт": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:         discordgo.MessageApplicationCommand,
				Name:         "Отправить репорт",
				DMPermission: new(bool),
			},
			Handler: reportMessageCommandHandler,
		},
		"update-guild": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:        discordgo.ChatApplicationCommand,
				Name:        "update-guild",
				Description: "Обновить настройки сервера",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "report-channel",
						Description: "Канал для репортов",
						ChannelTypes: []discordgo.ChannelType{
							discordgo.ChannelTypeGuildText,
						},
						Required: false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "resolved-report-channel",
						Description: "Канал для рассмотренных репортов",
						ChannelTypes: []discordgo.ChannelType{
							discordgo.ChannelTypeGuildText,
						},
						Required: false,
					},
				},
				DMPermission:             new(bool),
				DefaultMemberPermissions: &AdministratorPermission,
			},
			Handler: updateGuildChatCommandHandler,
		},
		"profile": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:        discordgo.ChatApplicationCommand,
				Name:        "profile",
				Description: "Просмотреть профиль пользователя",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "Пользователь",
						Required:    false,
					},
				},
				DMPermission: new(bool),
			},
			Handler: profileChatCommandHandler,
		},
		"+REP": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:         discordgo.UserApplicationCommand,
				Name:         "+REP",
				DMPermission: new(bool),
			},
			Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
				db.ChangeMemberReputation(interactionCreate.ApplicationCommandData().TargetID, 1)
			},
		},
		"-REP": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:         discordgo.UserApplicationCommand,
				Name:         "-REP",
				DMPermission: new(bool),
			},
			Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
				db.ChangeMemberReputation(interactionCreate.ApplicationCommandData().TargetID, -1)
			},
		},
	}
)

func CreateApplicationCommands(session *discordgo.Session) {
	for index, value := range Commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", value.ApplicationCommand)
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
	commands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		log.Printf("Error getting application commands: %v", err)
		return
	}
	for _, command := range commands {
		err = session.ApplicationCommandDelete(session.State.User.ID, "", command.ID)
		if err != nil {
			log.Panicf("Error deleting '%v' command: %v", command.Name, err)
		}
		log.Printf("Successfully deleted '%v' command", command.Name)
	}
}

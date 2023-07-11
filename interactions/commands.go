package interactions

import (
	"github.com/bwmarrin/discordgo"
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
		"guild": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:        discordgo.ChatApplicationCommand,
				Name:        "guild",
				Description: "Управление настройками сервера",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "update",
						Description: "Обновить настройки сервера",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionChannel,
								Name:        "канал_для_репортов",
								Description: "Канал, где находятся нерассмотренные репорты",
								ChannelTypes: []discordgo.ChannelType{
									discordgo.ChannelTypeGuildText,
								},
								Required: false,
							},
							{
								Type:        discordgo.ApplicationCommandOptionChannel,
								Name:        "канал_для_рассмотренных_репортов",
								Description: "Канал, где находятся рассмотренные репорты",
								ChannelTypes: []discordgo.ChannelType{
									discordgo.ChannelTypeGuildText,
								},
								Required: false,
							},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "view",
						Description: "Просмотреть настройки сервера",
					},
				},
				DMPermission:             new(bool),
				DefaultMemberPermissions: &AdministratorPermission,
			},
			Handler: guildChatCommandHandler,
		},
		"profile": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:        discordgo.ChatApplicationCommand,
				Name:        "profile",
				Description: "Просмотреть профиль пользователя",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "пользователь",
						Description: "Пользователь для просмотра профиля",
						Required:    false,
					},
				},
				DMPermission: new(bool),
			},
			Handler: profileChatCommandHandler,
		},
		"Лайк": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:         discordgo.UserApplicationCommand,
				Name:         "Лайк",
				DMPermission: new(bool),
			},
			Handler: likeUserCommandHandler,
		},
		"Дизлайк": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:         discordgo.UserApplicationCommand,
				Name:         "Дизлайк",
				DMPermission: new(bool),
			},
			Handler: dislikeUserCommandHandler,
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

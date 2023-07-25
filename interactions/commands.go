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
							{
								Type:        discordgo.ApplicationCommandOptionChannel,
								Name:        "канал_для_логирования_репутации",
								Description: "Канал, где логируется изменение репутации пользователей",
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
		"resetdelay": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:        discordgo.ChatApplicationCommand,
				Name:        "resetdelay",
				Description: "Сбросить задержку для лайков и дизлайков",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "пользователь",
						Description: "Пользователь, у которого вы хотите сбросить задержку",
						Required:    true,
					},
				},
				DMPermission:             new(bool),
				DefaultMemberPermissions: &AdministratorPermission,
			},
			Handler: resetDelayChatCommandHandler,
		},
		"Профиль": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:         discordgo.UserApplicationCommand,
				Name:         "Профиль",
				DMPermission: new(bool),
			},
			Handler: profileCommandHandler,
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
			Handler: profileCommandHandler,
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
		"like": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:        discordgo.ChatApplicationCommand,
				Name:        "like",
				Description: "Поставить лайк пользователю",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "пользователь",
						Description: "Пользователь, которому вы хотите поставить лайк",
						Required:    true,
					},
				},
				DMPermission: new(bool),
			},
			Handler: likeChatCommandHandler,
		},
		"dislike": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:        discordgo.ChatApplicationCommand,
				Name:        "dislike",
				Description: "Поставить дизлайк пользователю",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "пользователь",
						Description: "Пользователь, которому вы хотите поставить дизлайк",
						Required:    true,
					},
				},
				DMPermission: new(bool),
			},
			Handler: dislikeChatCommandHandler,
		},
		"top": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:        discordgo.ChatApplicationCommand,
				Name:        "top",
				Description: "Просмотреть топ пользователей",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "reputation",
						Description: "Просмотреть топ пользователей по репутации",
					},
				},
				DMPermission: new(bool),
			},
			Handler: topChatCommandHandler,
		},
		"setreputation": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:        discordgo.ChatApplicationCommand,
				Name:        "setreputation",
				Description: "Установить репутацию пользователю",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "пользователь",
						Description: "Пользователь, которому вы хотите установить репутацию",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						MinValue:    &db.MinReputation,
						MaxValue:    db.MaxReputation,
						Name:        "репутация",
						Description: "Репутация, которую вы хотите установить",
						Required:    true,
					},
				},
				DMPermission:             new(bool),
				DefaultMemberPermissions: &AdministratorPermission,
			},
			Handler: setReputationChatCommandHandler,
		},
		"warn": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:        discordgo.ChatApplicationCommand,
				Name:        "warn",
				Description: "Выдать предупреждение пользователю",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "пользователь",
						Description: "Пользователь, которому вы хотите выдать предупреждение",
						Required:    true,
					},
				},
			},
			Handler: warnChatCommandHandler,
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

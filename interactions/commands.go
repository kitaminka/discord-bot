package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"log"
)

type CommandHandler func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Type:         discordgo.MessageApplicationCommand,
			Name:         "Отправить репорт",
			DMPermission: new(bool),
		},
		GuildApplicationCommand,
		{
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
		{
			Type:         discordgo.UserApplicationCommand,
			Name:         "Профиль",
			DMPermission: new(bool),
		},
		{
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
		{
			Type:         discordgo.UserApplicationCommand,
			Name:         "Лайк",
			DMPermission: new(bool),
		},
		{
			Type:         discordgo.UserApplicationCommand,
			Name:         "Дизлайк",
			DMPermission: new(bool),
		},
		{
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
		{
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
		{
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
		{
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
		{
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
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "причина",
					Description: "Причина выдачи предупреждения",
					Required:    true,
					Choices:     reasonChoices,
				},
			},
			DMPermission:             new(bool),
			DefaultMemberPermissions: &ModeratorPermission,
		},
		{
			Type:                     discordgo.MessageApplicationCommand,
			Name:                     "Выдать предупреждение",
			DMPermission:             new(bool),
			DefaultMemberPermissions: &ModeratorPermission,
		},
		{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "remwarn",
			Description: "Снять предупреждения пользователя",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "пользователь",
					Description: "Пользователь, с которого вы хотите снять предупреждения",
					Required:    true,
				},
			},
			DMPermission:             new(bool),
			DefaultMemberPermissions: &ModeratorPermission,
		},
		{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "mute",
			Description: "Выдать мут пользователю",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "пользователь",
					Description: "Пользователь, которому вы хотите выдать мут",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "длительность",
					Description: "Длительность мута",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "причина",
					Description: "Причина выдачи мута",
					Choices:     reasonChoices,
					Required:    true,
				},
			},
			DMPermission:             new(bool),
			DefaultMemberPermissions: &ModeratorPermission,
		},
	}
	CommandHandlers = map[string]CommandHandler{
		"Отправить репорт": reportMessageCommandHandler,
		"guild":            guildChatCommandHandler,
		"resetdelay":       resetDelayChatCommandHandler,
		"Профиль":          profileCommandHandler,
		"profile":          profileCommandHandler,
		"Лайк":             likeUserCommandHandler,
		"Дизлайк":          dislikeUserCommandHandler,
		"like":             likeChatCommandHandler,
		"dislike":          dislikeChatCommandHandler,
		"top":              topChatCommandHandler,
		"setreputation":    setReputationChatCommandHandler,
		"warn":             warnChatCommandHandler,
		"remwarn":          remWarnChatCommandHandler,
		"Выдать предупреждение": warnMessageCommandHandler,
		"mute": muteChatCommandHandler,
	}
)

func CreateApplicationCommands(session *discordgo.Session) {
	commands, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, "", Commands)
	if err != nil {
		log.Panicf("Error creating commands: %v", err)
	}

	log.Print("Created commands")

	Commands = commands
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

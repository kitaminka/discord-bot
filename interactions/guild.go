package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/msg"
	"log"
)

var GuildApplicationCommand = &discordgo.ApplicationCommand{
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
}

func guildChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		InteractionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
		return
	}

	switch interactionCreate.ApplicationCommandData().Options[0].Name {
	case "view":
		guildViewChatCommandHandler(session, interactionCreate)
	case "update":
		guildUpdateChatCommandHandler(session, interactionCreate)
	default:
		InteractionRespondError(session, interactionCreate.Interaction, "Неизвестная подкоманда.")
	}
}
func guildViewChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}

	guild, err := db.GetGuild()
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при получении настроек сервера.")
		log.Printf("Error getting guild: %v", err)
		return
	}

	structuredDescriptionFields := []*msg.StructuredTextField{
		{
			Emoji: msg.ShieldCheckMarkEmoji,
			Name:  "ID сервера",
			Value: guild.ID,
		},
	}

	reportChannel, err := session.Channel(guild.ReportChannelID)
	if err != nil {
		log.Printf("Error getting report channel: %v", err)
	} else {
		structuredDescriptionFields = append(structuredDescriptionFields, &msg.StructuredTextField{
			Emoji: msg.ReportEmoji,
			Name:  "Канал для репортов",
			Value: reportChannel.Mention(),
		})
	}
	resolvedReportChannel, err := session.Channel(guild.ResoledReportChannelID)
	if err != nil {
		log.Printf("Error getting resolved report channel: %v", err)
	} else {
		structuredDescriptionFields = append(structuredDescriptionFields, &msg.StructuredTextField{
			Emoji: msg.ShieldCheckMarkEmoji,
			Name:  "Канал для рассмотренных репортов",
			Value: resolvedReportChannel.Mention(),
		})
	}
	reputationLogChannel, err := session.Channel(guild.ReputationLogChannelID)
	if err != nil {
		log.Printf("Error getting reputation log channel: %v", err)
	} else {
		structuredDescriptionFields = append(structuredDescriptionFields, &msg.StructuredTextField{
			Emoji: msg.ReputationEmoji,
			Name:  "Канал для логирования репутации",
			Value: reputationLogChannel.Mention(),
		})
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title: "Настройки сервера",
				Description: msg.StructuredText{
					Fields: structuredDescriptionFields,
				}.ToString(),
				Color: msg.DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Printf("Error editing interaction response: %v", err)
		return
	}
}
func guildUpdateChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if len(interactionCreate.ApplicationCommandData().Options[0].Options) == 0 {
		InteractionRespondError(session, interactionCreate.Interaction, "Вы не указали ни одной опции для обновления.")
		return
	}

	err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}

	structuredDescription := msg.StructuredText{
		Text: "Настройки сервера были успешно обновлены. Этот сервер установлен как основной.",
		Fields: []*msg.StructuredTextField{
			{
				Emoji: msg.IdEmoji,
				Name:  "ID сервера",
				Value: interactionCreate.GuildID,
			},
		},
	}
	server := db.Guild{
		ID: interactionCreate.GuildID,
	}

	for _, option := range interactionCreate.ApplicationCommandData().Options[0].Options {
		switch option.Name {
		case "канал_для_репортов":
			channel := option.ChannelValue(session)
			server.ReportChannelID = channel.ID
			structuredDescription.Fields = append(structuredDescription.Fields, &msg.StructuredTextField{
				Emoji: msg.ReportEmoji,
				Name:  "Канал для репортов",
				Value: channel.Mention(),
			})
		case "канал_для_рассмотренных_репортов":
			channel := option.ChannelValue(session)
			server.ResoledReportChannelID = channel.ID
			structuredDescription.Fields = append(structuredDescription.Fields, &msg.StructuredTextField{
				Emoji: msg.ShieldCheckMarkEmoji,
				Name:  "Канал для рассмотренных репортов",
				Value: channel.Mention(),
			})
		case "канал_для_логирования_репутации":
			channel := option.ChannelValue(session)
			server.ReputationLogChannelID = channel.ID
			structuredDescription.Fields = append(structuredDescription.Fields, &msg.StructuredTextField{
				Emoji: msg.ReputationEmoji,
				Name:  "Канал для логирования репутации",
				Value: channel.Mention(),
			})
		}
	}

	err = db.UpdateGuild(server)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при обновлении настроек сервера.")
		log.Printf("Error updating guild: %v", err)
		return
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "Настройки сервера обновлены",
				Description: structuredDescription.ToString(),
				Color:       msg.DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Printf("Error editing interaction response: %v", err)
		return
	}
}

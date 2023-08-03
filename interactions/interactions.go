package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/msg"
	"log"
	"strconv"
)

var (
	AdministratorPermission = int64(discordgo.PermissionAdministrator)
	ModeratorPermission     = int64(discordgo.PermissionModerateMembers)
)

func resetDelayChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		interactionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
		return
	}

	user := interactionCreate.ApplicationCommandData().Options[0].UserValue(session)
	if user.Bot {
		interactionRespondError(session, interactionCreate.Interaction, "Вы не можете сбросить задержку боту.")
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

	err = db.ResetUserReputationDelay(user.ID)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при сбросе задержки.")
		log.Printf("Error resetting user reputation delay: %v", err)
		return
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "Задержка сброшена",
				Description: fmt.Sprintf("Задержка пользователя %v была сброшена.", msg.UserMention(user)),
				Color:       msg.DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Printf("Error editing interaction response: %v", err)
		return
	}
}
func guildChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		interactionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
		return
	}

	switch interactionCreate.ApplicationCommandData().Options[0].Name {
	case "view":
		guildViewChatCommandHandler(session, interactionCreate)
	case "update":
		guildUpdateChatCommandHandler(session, interactionCreate)
	case "rules":
		guildRulesChatCommandHandler(session, interactionCreate)
	default:
		interactionRespondError(session, interactionCreate.Interaction, "Неизвестная подкоманда.")
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

	structuredDescriptionFields := []*msg.StructuredDescriptionField{
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
		structuredDescriptionFields = append(structuredDescriptionFields, &msg.StructuredDescriptionField{
			Emoji: msg.ReportEmoji,
			Name:  "Канал для репортов",
			Value: reportChannel.Mention(),
		})
	}
	resolvedReportChannel, err := session.Channel(guild.ResoledReportChannelID)
	if err != nil {
		log.Printf("Error getting resolved report channel: %v", err)
	} else {
		structuredDescriptionFields = append(structuredDescriptionFields, &msg.StructuredDescriptionField{
			Emoji: msg.ShieldCheckMarkEmoji,
			Name:  "Канал для рассмотренных репортов",
			Value: resolvedReportChannel.Mention(),
		})
	}
	reputationLogChannel, err := session.Channel(guild.ReputationLogChannelID)
	if err != nil {
		log.Printf("Error getting reputation log channel: %v", err)
	} else {
		structuredDescriptionFields = append(structuredDescriptionFields, &msg.StructuredDescriptionField{
			Emoji: msg.ReputationEmoji,
			Name:  "Канал для логирования репутации",
			Value: reputationLogChannel.Mention(),
		})
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title: "Настройки сервера",
				Description: msg.StructuredDescription{
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
		interactionRespondError(session, interactionCreate.Interaction, "Вы не указали ни одной опции для обновления.")
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

	structuredDescription := msg.StructuredDescription{
		Text: "Настройки сервера были успешно обновлены. Этот сервер установлен как основной.",
		Fields: []*msg.StructuredDescriptionField{
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
			structuredDescription.Fields = append(structuredDescription.Fields, &msg.StructuredDescriptionField{
				Emoji: msg.ReportEmoji,
				Name:  "Канал для репортов",
				Value: channel.Mention(),
			})
		case "канал_для_рассмотренных_репортов":
			channel := option.ChannelValue(session)
			server.ResoledReportChannelID = channel.ID
			structuredDescription.Fields = append(structuredDescription.Fields, &msg.StructuredDescriptionField{
				Emoji: msg.ShieldCheckMarkEmoji,
				Name:  "Канал для рассмотренных репортов",
				Value: channel.Mention(),
			})
		case "канал_для_логирования_репутации":
			channel := option.ChannelValue(session)
			server.ReputationLogChannelID = channel.ID
			structuredDescription.Fields = append(structuredDescription.Fields, &msg.StructuredDescriptionField{
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
func guildRulesChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	switch interactionCreate.ApplicationCommandData().Options[0].Options[0].Name {
	case "add":
		guildRulesAddChatCommandHandler(session, interactionCreate)
	case "view":
		guildRulesViewChatCommandHandler(session, interactionCreate)
	}
}
func guildRulesAddChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
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

	var name, description string

	for _, option := range interactionCreate.ApplicationCommandData().Options[0].Options[0].Options {
		switch option.Name {
		case "название":
			name = option.StringValue()
		case "описание":
			description = option.StringValue()
		}
	}

	err = db.AddGuildRule(db.Rule{
		Name:        name,
		Description: description,
	})
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при добавлении правила.")
		log.Printf("Error adding guild rule: %v", err)
		return
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title: "Правило добавлено",
				Description: msg.StructuredDescription{
					Text: "Правило было успешно добавлено.",
					Fields: []*msg.StructuredDescriptionField{
						{
							Name:  "Название",
							Value: name,
						},
						{
							Name:  "Описание",
							Value: description,
						},
					},
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
func guildRulesViewChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {

}

// Used for user command and chat command
func profileCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
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

	var member *discordgo.Member

	if len(interactionCreate.ApplicationCommandData().TargetID) != 0 {
		member, err = session.GuildMember(interactionCreate.GuildID, interactionCreate.ApplicationCommandData().TargetID)
		if err != nil {
			interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при получении профиля пользователя. Свяжитесь с администрацией.")
			log.Printf("Error getting member: %v", err)
			return
		}
	} else if len(interactionCreate.ApplicationCommandData().Options) == 0 {
		member = interactionCreate.Member
	} else {
		discordUser := interactionCreate.ApplicationCommandData().Options[0].UserValue(session)

		member, err = session.GuildMember(interactionCreate.GuildID, discordUser.ID)
		if err != nil {
			interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при получении профиля пользователя. Свяжитесь с администрацией.")
			log.Printf("Error getting member: %v", err)
			return
		}
	}

	if member.User.Bot {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Вы не можете просмотреть профиль бота.")
		return
	}

	user, err := db.GetUser(member.User.ID)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при получении профиля пользователя. Свяжитесь с администрацией.")
		log.Printf("Error getting user: %v", err)
		return
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title: fmt.Sprintf("%v Профиль пользователя %v", msg.UserEmoji, member.User.Username),
				Description: msg.StructuredDescription{
					Fields: []*msg.StructuredDescriptionField{
						{
							Emoji: msg.UsernameEmoji,
							Name:  "Пользователь",
							Value: member.Mention(),
						},
						{
							Emoji: msg.JoinEmoji,
							Name:  "Присоединился к серверу",
							Value: fmt.Sprintf("<t:%v:R>", member.JoinedAt.Unix()),
						},
						{
							Emoji: msg.ReputationEmoji,
							Name:  "Репутация",
							Value: strconv.Itoa(user.Reputation),
						},
						{
							Emoji: msg.ReportEmoji,
							Name:  "Отправленные репортов",
							Value: strconv.Itoa(user.ReportsSentCount),
						},
					},
				}.ToString(),
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: member.AvatarURL(""),
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("ID: %v", member.User.ID),
				},
				Color: msg.DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Printf("Error editing interaction response: %v", err)
		return
	}
}

func topChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	switch interactionCreate.ApplicationCommandData().Options[0].Name {
	case "reputation":
		topReputationChatCommandHandler(session, interactionCreate)
	default:
		interactionRespondError(session, interactionCreate.Interaction, "Неизвестная подкоманда.")
	}
}

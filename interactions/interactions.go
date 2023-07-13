package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/msg"
	"log"
	"strconv"
)

var AdministratorPermission = int64(discordgo.PermissionAdministrator)

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
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при сбросе задержки.")
		log.Printf("Error resetting user reputation delay: %v", err)
		return
	}

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Задержка сброшена",
				Description: fmt.Sprintf("Задержка пользователя %v была сброшена.", msg.UserMention(user)),
				Color:       msg.DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Printf("Error creating followup message: %v", err)
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

	guild, err := db.GetGuild(interactionCreate.GuildID)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при получении настроек сервера.")
		log.Printf("Error getting guild: %v", err)
		return
	}

	reportChannel, err := session.Channel(guild.ReportChannelID)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при получении настроек сервера.")
		log.Printf("Error getting report channel: %v", err)
		return
	}
	resolvedReportChannel, err := session.Channel(guild.ResoledReportChannelID)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при получении настроек сервера.")
		log.Printf("Error getting resolved report channel: %v", err)
		return
	}

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Настройки сервера",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "ID сервера",
						Value: guild.ID,
					},
					{
						Name:  "Канал для репортов",
						Value: reportChannel.Mention(),
					},
					{
						Name:  "Канал для рассмотренные репортов",
						Value: resolvedReportChannel.Mention(),
					},
				},
				Color: msg.DefaultEmbedColor,
			},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		log.Printf("Error creating followup message: %v", err)
		return
	}
}
func guildUpdateChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
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

	var fields []*discordgo.MessageEmbedField
	server := db.Guild{
		ID: interactionCreate.GuildID,
	}

	for _, option := range interactionCreate.ApplicationCommandData().Options[0].Options {
		switch option.Name {
		case "канал_для_репортов":
			channel := option.ChannelValue(session)
			server.ReportChannelID = channel.ID
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  "Канал для репортов",
				Value: channel.Mention(),
			})
		case "канал_для_рассмотренных_репортов":
			channel := option.ChannelValue(session)
			server.ResoledReportChannelID = channel.ID
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  "Канал для рассмотренных репортов",
				Value: channel.Mention(),
			})
		case "канал_для_логирования_репутации":
			channel := option.ChannelValue(session)
			server.ReputationLogChannelID = channel.ID
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  "Канал для логирования репутации",
				Value: channel.Mention(),
			})
		}
	}

	err = db.UpdateGuild(server)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при обновлении настроек сервера.")
		log.Printf("Error updating guild: %v", err)
		return
	}

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Настройки сервера обновлены",
				Description: "Настройки сервера были успешно обновлены.",
				Fields:      fields,
				Color:       msg.DefaultEmbedColor,
			},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		log.Printf("Error creating followup message: %v", err)
		return
	}
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
			followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при получении профиля пользователя. Свяжитесь с администрацией.")
			log.Printf("Error getting member: %v", err)
			return
		}
	} else if len(interactionCreate.ApplicationCommandData().Options) == 0 {
		member = interactionCreate.Member
	} else {
		discordUser := interactionCreate.ApplicationCommandData().Options[0].UserValue(session)

		member, err = session.GuildMember(interactionCreate.GuildID, discordUser.ID)
		if err != nil {
			followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при получении профиля пользователя. Свяжитесь с администрацией.")
			log.Printf("Error getting member: %v", err)
			return
		}
	}

	if member.User.Bot {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Вы не можете просмотреть профиль бота.")
		return
	}

	user, err := db.GetUser(member.User.ID)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при получении профиля пользователя. Свяжитесь с администрацией.")
		log.Printf("Error getting user: %v", err)
		return
	}

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       fmt.Sprintf("%v Профиль пользователя %v", msg.UserEmoji, member.User.Username),
				Description: fmt.Sprintf("%v **Пользователь**: %v\n%v **Присоединился к серверу**: <t:%v:R>\n%v **Репутация**: %v\n%v **Отправленные репортов**: %v\n", msg.MentionEmoji, member.Mention(), msg.JoinEmoji, member.JoinedAt.Unix(), msg.ReputationEmoji, user.Reputation, msg.ReportEmoji, user.ReportsSentCount),
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
		log.Printf("Error creating followup message: %v", err)
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

func setReputationChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		interactionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
		return
	}

	var targetUser *discordgo.User
	var reputation int

	for _, option := range interactionCreate.ApplicationCommandData().Options {
		switch option.Name {
		case "пользователь":
			targetUser = option.UserValue(session)
		case "репутация":
			reputation = int(option.IntValue())
		}
	}

	if targetUser == nil {
		interactionRespondError(session, interactionCreate.Interaction, "Не указан пользователь.")
		return
	} else if reputation == 0 {
		interactionRespondError(session, interactionCreate.Interaction, "Не указана репутация.")
		return
	}

	if targetUser.Bot {
		interactionRespondError(session, interactionCreate.Interaction, "Вы не можете изменить репутацию бота.")
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

	err = db.SetUserReputation(targetUser.ID, reputation)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при изменении репутации пользователя. Свяжитесь с администрацией.")
		log.Printf("Error setting user reputation: %v", err)
		return
	}

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Репутация изменена",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Пользователь",
						Value: msg.UserMention(targetUser),
					},
					{
						Name:  "Репутация",
						Value: strconv.Itoa(reputation),
					},
				},
				Color: msg.DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Printf("Error creating followup message: %v", err)
		return
	}
}

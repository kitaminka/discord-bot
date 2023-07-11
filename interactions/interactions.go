package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"log"
	"strconv"
)

const (
	DefaultEmbedColor = 14546431
	ErrorEmbedColor   = 16711680
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
				Description: fmt.Sprintf("Задержка пользователя %v была сброшена.", user.Mention()),
				Color:       DefaultEmbedColor,
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
				Color: DefaultEmbedColor,
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

	for _, option := range interactionCreate.ApplicationCommandData().Options {
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
				Color:       DefaultEmbedColor,
			},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		log.Printf("Error creating followup message: %v", err)
		return
	}
}
func profileChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}

	var user *discordgo.User

	if len(interactionCreate.ApplicationCommandData().Options) == 0 {
		user = interactionCreate.Member.User
	} else {
		user = interactionCreate.ApplicationCommandData().Options[0].UserValue(session)
	}

	if user.Bot {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Вы не можете просмотреть профиль бота.")
		return
	}

	member, err := db.GetUser(user.ID)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при получении профиля пользователя.")
		log.Printf("Error getting member: %v", err)
		return
	}

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Профиль пользователя",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Пользователь",
						Value: user.Mention(),
					},
					{
						Name:  "ID пользователя",
						Value: user.ID,
					},
					{
						Name:  "Репутация",
						Value: strconv.Itoa(member.Reputation),
					},
					{
						Name:  "Количство отправленных репортов",
						Value: strconv.Itoa(member.ReportsSentCount),
					},
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: user.AvatarURL(""),
				},
				Color: DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Printf("Error creating followup message: %v", err)
		return
	}
}

package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/cfg"
	"github.com/kitaminka/discord-bot/db"
	"log"
)

func reportMessageCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("Error responding to interaction: ", err)
		return
	}

	guild, err := db.GetGuild(interactionCreate.GuildID)
	if err != nil {
		log.Println("Error getting server: ", err)
		return
	}

	reportedMessage := interactionCreate.ApplicationCommandData().Resolved.Messages[interactionCreate.ApplicationCommandData().TargetID]
	reportedMessageContent := fmt.Sprintf("```%v```", reportedMessage.Content)
	reportedMessageUrl := fmt.Sprintf("https://discord.com/channels/%v/%v/%v", interactionCreate.GuildID, interactionCreate.ChannelID, reportedMessage.ID)
	reportedMessageSenderMention := fmt.Sprintf("<@%v>", reportedMessage.Author.ID)
	reportSenderMention := fmt.Sprintf("<@%v>", interactionCreate.Member.User.ID)

	_, err = session.ChannelMessageSendComplex(guild.ReportChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Репорт",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Отправитель репорта",
						Value: reportSenderMention,
					},
					{
						Name:  "Сообщение",
						Value: reportedMessageUrl,
					},
					{
						Name:  "Отправитель сообщения",
						Value: reportedMessageSenderMention,
					},
					{
						Name:  "Содержимое сообщения",
						Value: reportedMessageContent,
					},
				},
				Color: DefaultEmbedColor,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					ResolveReportButton,
				},
			},
		},
	})
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при отправке вашего репорта. Свяжитесь с администрацией.")
		log.Println("Error sending report: ", err)
		return
	}

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Репорт отправлен",
				Description: "Ваш репорт был успешно отправлен.",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Сообщение",
						Value: reportedMessageUrl,
					},
					{
						Name:  "Отправитель сообщения",
						Value: reportedMessageSenderMention,
					},
				},
				Color: DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Println("Error creating followup message: ", err)
		return
	}
}

func resolveReportHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("Error responding to interaction: ", err)
		return
	}

	if len(interactionCreate.Message.Embeds) != 1 {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при рассмотрении репорта. Свяжитесь с администрацией.")
		log.Println("Report message is invalid")
		return
	}

	reportResolverMention := fmt.Sprintf("<@%v>", interactionCreate.Member.User.ID)
	reportMessageEmbed := interactionCreate.Message.Embeds[0]

	resolvedReportMessage, err := session.ChannelMessageSendComplex(cfg.Config.ResoledReportChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Рассмотренный репорт",
				Fields: append(reportMessageEmbed.Fields, &discordgo.MessageEmbedField{
					Name:  "Рассмотритель",
					Value: reportResolverMention,
				}),
				Color: DefaultEmbedColor,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					ReturnReportButton,
				},
			},
		},
	})
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при рассмотрении репорта. Свяжитесь с администрацией.")
		log.Println("Error sending resolved report: ", err)
		return
	}
	err = session.ChannelMessageDelete(interactionCreate.Message.ChannelID, interactionCreate.Message.ID)
	if err != nil {
		err = session.ChannelMessageDelete(resolvedReportMessage.ChannelID, resolvedReportMessage.ID)
		if err != nil {
			log.Println("Error deleting resolved report: ", err)
		}
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при рассмотрении репорта. Свяжитесь с администрацией.")
		log.Println("Error deleting report: ", err)
		return
	}

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Репорт рассмотрен",
				Description: "Репорт был успешно перемещен в рассмотренные.",
				Color:       DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Println("Error creating followup message: ", err)
		return
	}
}

func returnReportHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("Error responding to interaction: ", err)
		return
	}

	if len(interactionCreate.Message.Embeds) != 1 {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при возвращении репорта. Свяжитесь с администрацией.")
		log.Println("Resolved report message is invalid")
		return
	}

	resolvedReportMessageEmbed := interactionCreate.Message.Embeds[0]

	resolvedReportMessage, err := session.ChannelMessageSendComplex(cfg.Config.ReportChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:  "Репорт",
				Fields: resolvedReportMessageEmbed.Fields[:len(resolvedReportMessageEmbed.Fields)-1],
				Color:  DefaultEmbedColor,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					ResolveReportButton,
				},
			},
		},
	})
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при возвращении репорта. Свяжитесь с администрацией.")
		log.Println("Error sending report: ", err)
		return
	}
	err = session.ChannelMessageDelete(interactionCreate.Message.ChannelID, interactionCreate.Message.ID)
	if err != nil {
		err = session.ChannelMessageDelete(resolvedReportMessage.ChannelID, resolvedReportMessage.ID)
		if err != nil {
			log.Println("Error deleting report: ", err)
		}
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при возвращении репорта. Свяжитесь с администрацией.")
		log.Println("Error deleting resolved report: ", err)
		return
	}

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Репорт возвращен",
				Description: "Репорт был успешно возвращен в нерассмотренные.",
				Color:       DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Println("Error creating followup message: ", err)
		return
	}
}

package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/cfg"
	"log"
)

func reportMessageCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	reportedMessage := interactionCreate.ApplicationCommandData().Resolved.Messages[interactionCreate.ApplicationCommandData().TargetID]
	reportedMessageContent := fmt.Sprintf("```%v```", reportedMessage.Content)
	reportedMessageUrl := fmt.Sprintf("https://discord.com/channels/%v/%v/%v", interactionCreate.GuildID, interactionCreate.ChannelID, reportedMessage.ID)
	reportedMessageSenderMention := fmt.Sprintf("<@%v>", reportedMessage.Author.ID)
	reportSenderMention := fmt.Sprintf("<@%v>", interactionCreate.Member.User.ID)

	err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("Error responding to interaction: ", err)
	}

	_, err = session.ChannelMessageSendComplex("1121453163451514880", &discordgo.MessageSend{
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
				Color: cfg.Config.EmbedColors.Default,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					Components["resolve_report"].MessageComponent,
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
				Color: cfg.Config.EmbedColors.Default,
			},
		},
	})
	if err != nil {
		log.Println("Error creating followup message: ", err)
	}
}

func resolveReportHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if len(interactionCreate.Message.Embeds) != 1 {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при рассмотрении репорта. Свяжитесь с администрацией.")
		log.Println("Report message is invalid")
		return
	}

	reportResolverMention := fmt.Sprintf("<@%v>", interactionCreate.Member.User.ID)
	reportMessageEmbed := interactionCreate.Message.Embeds[0]

	err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("Error responding to interaction: ", err)
	}

	resolvedReportMessage, err := session.ChannelMessageSendComplex("1122193280445194340", &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Рассмотренный репорт",
				Fields: append(reportMessageEmbed.Fields, &discordgo.MessageEmbedField{
					Name:  "Рассмотритель",
					Value: reportResolverMention,
				}),
				Color: cfg.Config.EmbedColors.Default,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					&discordgo.Button{
						Label:    "Вернуть",
						Style:    discordgo.PrimaryButton,
						CustomID: "return_report",
						Emoji: discordgo.ComponentEmoji{
							Name: "🔄",
						},
					},
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
				Color:       cfg.Config.EmbedColors.Default,
			},
		},
	})
	if err != nil {
		log.Println("Error creating followup message: ", err)
	}
}

func returnReportHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if len(interactionCreate.Message.Embeds) != 1 {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при возвращении репорта. Свяжитесь с администрацией.")
		log.Println("Resolved report message is invalid")
		return
	}

	resolvedReportMessageEmbed := interactionCreate.Message.Embeds[0]

	err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("Error responding to interaction: ", err)
	}

	resolvedReportMessage, err := session.ChannelMessageSendComplex("1121453163451514880", &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:  "Репорт",
				Fields: resolvedReportMessageEmbed.Fields[:len(resolvedReportMessageEmbed.Fields)-1],
				Color:  cfg.Config.EmbedColors.Default,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					&discordgo.Button{
						CustomID: "resolve_report",
						Label:    "Рассмотрено",
						Style:    discordgo.SuccessButton,
						Emoji: discordgo.ComponentEmoji{
							Name: "✅",
						},
					},
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
				Color:       cfg.Config.EmbedColors.Default,
			},
		},
	})
	if err != nil {
		log.Println("Error creating followup message: ", err)
	}
}

package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func reportMessageCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	targetMessage := interactionCreate.ApplicationCommandData().Resolved.Messages[interactionCreate.ApplicationCommandData().TargetID]
	targetMessageContent := fmt.Sprintf("```%v```", targetMessage.Content)
	targetMessageUrl := fmt.Sprintf("https://discord.com/channels/%v/%v/%v", interactionCreate.GuildID, interactionCreate.ChannelID, targetMessage.ID)
	targetMessageSenderMention := fmt.Sprintf("<@%v>", targetMessage.Author.ID)
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
				Title: "Новый репорт",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Отправитель репорта",
						Value: reportSenderMention,
					},
					{
						Name:  "Сообщение",
						Value: targetMessageUrl,
					},
					{
						Name:  "Отправитель сообщения",
						Value: targetMessageSenderMention,
					},
					{
						Name:  "Содержимое сообщения",
						Value: targetMessageContent,
					},
				},
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Рассмотрено",
						Style:    discordgo.SuccessButton,
						CustomID: "report_resolved",
						Emoji: discordgo.ComponentEmoji{
							Name: "✅",
						},
					},
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
						Value: targetMessageUrl,
					},
					{
						Name:  "Отправитель сообщения",
						Value: targetMessageSenderMention,
					},
				},
			},
		},
	})
	if err != nil {
		log.Println("Error creating followup message: ", err)
	}
}

func reportResolvedHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if len(interactionCreate.Message.Embeds) != 1 {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при рассмотрении репорта. Свяжитесь с администрацией.")
		log.Println("Report message is invalid")
		return
	}

	reportResolverMention := fmt.Sprintf("<@%v>", interactionCreate.Member.User.ID)
	reportMessageEmbeds := interactionCreate.Message.Embeds[0]

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
				Fields: append(reportMessageEmbeds.Fields, &discordgo.MessageEmbedField{
					Name:  "Рассмотритель",
					Value: reportResolverMention,
				}),
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
			},
		},
	})
	if err != nil {
		log.Println("Error creating followup message: ", err)
	}
}

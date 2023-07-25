package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/msg"
	"log"
)

func reportMessageCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	reportedMessage := interactionCreate.ApplicationCommandData().Resolved.Messages[interactionCreate.ApplicationCommandData().TargetID]
	reportedMessageContent := fmt.Sprintf("```%v```", reportedMessage.Content)
	reportedMessageUrl := fmt.Sprintf("https://discord.com/channels/%v/%v/%v", interactionCreate.GuildID, interactionCreate.ChannelID, reportedMessage.ID)
	reportedMessageSenderMention := msg.UserMention(reportedMessage.Author)
	reportSenderMention := msg.UserMention(interactionCreate.Member.User)

	if interactionCreate.Member.User.ID == reportedMessage.Author.ID {
		interactionRespondError(session, interactionCreate.Interaction, "Вы не можете отправить репорт на своё сообщение.")
		return
	}

	if reportedMessage.Author.Bot {
		interactionRespondError(session, interactionCreate.Interaction, "Вы не можете отправить репорт на сообщение бота.")
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

	guild, err := db.GetGuild()
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при отправке вашего репорта. Свяжитесь с администрацией.")
		log.Printf("Error getting server: %v", err)
		return
	}

	_, err = session.ChannelMessageSendComplex(guild.ReportChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: fmt.Sprintf("%v Новый репорт", msg.ReportEmoji),
				Description: msg.StructuredDescription{
					Fields: []*msg.StructuredDescriptionField{
						{
							Emoji: msg.UsernameEmoji,
							Name:  "Отправитель репорта",
							Value: reportSenderMention,
						},
						{
							Emoji: msg.TextChannelEmoji,
							Name:  "Сообщение",
							Value: reportedMessageUrl,
						},
						{
							Emoji: msg.UserEmoji,
							Name:  "Отправитель сообщения",
							Value: reportedMessageSenderMention,
						},
						{
							Emoji: msg.MessageEmoji,
							Name:  "Содержимое сообщения",
							Value: reportedMessageContent,
						},
					},
				}.ToString(),
				Color: msg.DefaultEmbedColor,
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
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при отправке вашего репорта. Свяжитесь с администрацией.")
		log.Printf("Error sending report: %v", err)
		return
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title: fmt.Sprintf("%v Репорт отправлен", msg.ReportEmoji),
				Description: msg.StructuredDescription{
					Text: "Ваш репорт был успешно отправлен.",
					Fields: []*msg.StructuredDescriptionField{
						{
							Emoji: msg.TextChannelEmoji,
							Name:  "Сообщение",
							Value: reportedMessageUrl,
						},
						{
							Emoji: msg.UserEmoji,
							Name:  "Отправитель сообщения",
							Value: reportedMessageSenderMention,
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
	err = db.IncrementUserReportsSent(interactionCreate.Member.User.ID)
	if err != nil {
		log.Printf("Error incrementing user reports sent: %v", err)
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
		log.Printf("Error responding to interaction: %v", err)
		return
	}

	guild, err := db.GetGuild()
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при рассмотрении репорта. Свяжитесь с администрацией.")
		log.Printf("Error getting server: %v", err)
		return
	}

	if len(interactionCreate.Message.Embeds) != 1 {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при рассмотрении репорта. Свяжитесь с администрацией.")
		log.Print("Report message is invalid")
		return
	}

	reportResolverMember := interactionCreate.Member
	reportMessageEmbed := interactionCreate.Message.Embeds[0]

	resolvedReportMessage, err := session.ChannelMessageSendComplex(guild.ResoledReportChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       fmt.Sprintf("%v Рассмотрен репорт", msg.ShieldCheckMarkEmoji),
				Description: reportMessageEmbed.Description,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    fmt.Sprintf("Рассмотрено: %v", reportResolverMember.User.Username),
					IconURL: reportResolverMember.AvatarURL(""),
				},
				Color: msg.DefaultEmbedColor,
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
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при рассмотрении репорта. Свяжитесь с администрацией.")
		log.Printf("Error sending resolved report: %v", err)
		return
	}
	err = session.ChannelMessageDelete(interactionCreate.Message.ChannelID, interactionCreate.Message.ID)
	if err != nil {
		err = session.ChannelMessageDelete(resolvedReportMessage.ChannelID, resolvedReportMessage.ID)
		if err != nil {
			log.Printf("Error deleting resolved report: %v", err)
		}
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при рассмотрении репорта. Свяжитесь с администрацией.")
		log.Printf("Error deleting report: %v", err)
		return
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "Репорт рассмотрен",
				Description: "Репорт был успешно перемещен в рассмотренные.",
				Color:       msg.DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Printf("Error editing interaction response: %v", err)
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
		log.Printf("Error responding to interaction: %v", err)
		return
	}

	guild, err := db.GetGuild()
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при возвращении репорта. Свяжитесь с администрацией.")
		log.Printf("Error getting server: %v", err)
		return
	}

	if len(interactionCreate.Message.Embeds) != 1 {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при возвращении репорта. Свяжитесь с администрацией.")
		log.Print("Resolved report message is invalid")
		return
	}

	resolvedReportMessageEmbed := interactionCreate.Message.Embeds[0]

	resolvedReportMessage, err := session.ChannelMessageSendComplex(guild.ReportChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Репорт",
				Description: resolvedReportMessageEmbed.Description,
				Color:       msg.DefaultEmbedColor,
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
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при возвращении репорта. Свяжитесь с администрацией.")
		log.Printf("Error sending report: %v", err)
		return
	}
	err = session.ChannelMessageDelete(interactionCreate.Message.ChannelID, interactionCreate.Message.ID)
	if err != nil {
		err = session.ChannelMessageDelete(resolvedReportMessage.ChannelID, resolvedReportMessage.ID)
		if err != nil {
			log.Printf("Error deleting report: %v", err)
		}
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при возвращении репорта. Свяжитесь с администрацией.")
		log.Printf("Error deleting resolved report: %v", err)
		return
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "Репорт возвращен",
				Description: "Репорт был успешно возвращен в нерассмотренные.",
				Color:       msg.DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Printf("Error editing interaction response: %v", err)
		return
	}
}

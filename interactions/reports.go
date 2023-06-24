package interactions

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func reportMessageCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	targetMessage := interactionCreate.ApplicationCommandData().Resolved.Messages[interactionCreate.ApplicationCommandData().TargetID]
	targetMessageUrl := "https://discord.com/channels/" + interactionCreate.GuildID + "/" + interactionCreate.ChannelID + "/" + targetMessage.ID
	targetMessageSenderMention := "<@" + targetMessage.Author.ID + ">"
	reportSenderMention := "<@" + interactionCreate.Member.User.ID + ">"

	err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("Error responding to interaction: ", err)
	}

	session.ChannelMessageSendComplex("1121453163451514880", &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Новый репорт",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Сообщение",
						Value: targetMessageUrl,
					},
					{
						Name:  "Отправитель сообщения",
						Value: targetMessageSenderMention,
					},
					{
						Name:  "Отправитель репорта",
						Value: reportSenderMention,
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

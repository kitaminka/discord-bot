package interactions

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type Component struct {
	MessageComponent discordgo.MessageComponent
	Handler          func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
}

var Components = map[string]Component{
	"report_resolved": {
		MessageComponent: &discordgo.Button{
			CustomID: "report_resolved",
		},
		Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
			reportMessageEmbeds := interactionCreate.Message.Embeds

			err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				log.Println("Error responding to interaction: ", err)
			}

			session.ChannelMessageSendComplex("1122193280445194340", &discordgo.MessageSend{
				Embeds: reportMessageEmbeds,
			})
			session.ChannelMessageDelete(interactionCreate.Message.ChannelID, interactionCreate.Message.ID)

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
		},
	},
}

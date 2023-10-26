package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/msg"
	"log"
)

func infoChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Информация о боте",
					Description: msg.StructuredText{
						Fields: []*msg.StructuredTextField{
							{
								Name:  "Разработчик",
								Value: "<@890320305082478652>",
							},
							{
								Name:  "Версия DiscordGo",
								Value: discordgo.VERSION,
							},
						},
					}.ToString(),
					Color: msg.DefaultEmbedColor,
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}
}

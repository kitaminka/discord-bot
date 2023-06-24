package interactions

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func createErrorEmbed(errorMessage string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Ошибка",
		Description: errorMessage,
	}
}
func interactionRespondError(session *discordgo.Session, interaction *discordgo.Interaction, errorMessage string) {
	err := session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				createErrorEmbed(errorMessage),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}
func followupErrorMessageCreate(session *discordgo.Session, interaction *discordgo.Interaction, errorMessage string) {
	_, err := session.FollowupMessageCreate(interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			createErrorEmbed(errorMessage),
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

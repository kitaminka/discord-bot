package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/msg"
	"log"
)

func createErrorEmbed(errorMessage string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%v Ошибка", msg.ErrorEmoji),
		Description: errorMessage,
		Color:       msg.ErrorEmbedColor,
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
func interactionResponseErrorEdit(session *discordgo.Session, interaction *discordgo.Interaction, errorMessage string) {
	_, err := session.InteractionResponseEdit(interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			createErrorEmbed(errorMessage),
		},
	})
	if err != nil {
		log.Printf("Error editing interaction response: %v", err)
	}
}
func followupErrorMessageCreate(session *discordgo.Session, interaction *discordgo.Interaction, errorMessage string) {
	_, err := session.InteractionResponseEdit(interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			createErrorEmbed(errorMessage),
		},
	})
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

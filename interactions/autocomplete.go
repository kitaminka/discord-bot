package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"log"
)

var AutocompleteHandlers = map[string]func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate){
	"warn": autocompleteWarnHandler,
}

func autocompleteWarnHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	var reasons []db.Reason
	var reasonChoices []*discordgo.ApplicationCommandOptionChoice

	reasons, err := db.GetReasons()
	if err != nil {
		log.Printf("Error getting guild: %v", err)
		return
	}

	for _, reason := range reasons {
		reasonChoices = append(reasonChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  reason.Name,
			Value: reason.ID.Hex(),
		})
	}

	err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: reasonChoices,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}
}

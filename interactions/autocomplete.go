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
	guild, err := db.GetGuild()
	if err != nil {
		log.Printf("Error getting guild: %v", err)
		return
	}

	var ruleChoices []*discordgo.ApplicationCommandOptionChoice

	for _, rule := range guild.Rules {
		ruleChoices = append(ruleChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  rule.Name,
			Value: rule.ID.Hex(),
		})
	}

	err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: ruleChoices,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}
}

package setup

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/interactions"
)

func SetupCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	app, err := session.Application("@me")
	if err != nil {
		interactions.InteractionRespondError(session, interactionCreate.Interaction, "Произошла ошибка.")
		log.Printf("Error getting application: %v", err)
		return
	}

	if interactionCreate.Member.User.ID != app.Owner.ID {
		interactions.InteractionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
		return
	}

	interactions.CreateApplicationCommands(session, interactionCreate.GuildID)

	session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Meow :3",
		},
	})

	time.Sleep(10 * time.Second)

	interactions.DeleteApplicationCommands(session, interactionCreate.GuildID)
}

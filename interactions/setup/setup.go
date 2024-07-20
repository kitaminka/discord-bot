package setup

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/interactions"
)

var SetupCommand = &discordgo.ApplicationCommand{
	Type:                     discordgo.ChatApplicationCommand,
	Name:                     "setup",
	Description:              "Начальная настройка бота",
	DMPermission:             new(bool),
	DefaultMemberPermissions: &interactions.AdministratorPermission,
}

func CreateSetupCommand(session *discordgo.Session) {
	_, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, "", []*discordgo.ApplicationCommand{SetupCommand})
	if err != nil {
		log.Panicf("Error creating setup command: %v", err)
	}

	log.Print("Setup command created")
}

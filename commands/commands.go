package commands

import "github.com/bwmarrin/discordgo"

type CommandHandler func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)

type Command struct {
	ApplicationCommand *discordgo.ApplicationCommand
	Handler            CommandHandler
}

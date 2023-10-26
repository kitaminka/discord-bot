package interactions

import (
	"github.com/bwmarrin/discordgo"
)

var AutocompleteHandlers = map[string]func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate){}

package interactions

import (
	"github.com/bwmarrin/discordgo"
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
		Handler: reportResolvedHandler,
	},
}

package interactions

import (
	"github.com/bwmarrin/discordgo"
)

type Component struct {
	MessageComponent discordgo.MessageComponent
	Handler          func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
}

var Components = map[string]Component{
	"resolve_report": {
		MessageComponent: &discordgo.Button{
			CustomID: "report_resolved",
		},
		Handler: resolveReportHandler,
	},
	"return_report": {
		MessageComponent: &discordgo.Button{
			CustomID: "report_resolved",
		},
		Handler: returnReportHandler,
	},
}

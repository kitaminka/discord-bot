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
			CustomID: "resolve_report",
			Label:    "Рассмотрено",
			Style:    discordgo.SuccessButton,
			Emoji: discordgo.ComponentEmoji{
				Name: "✅",
			},
		},
		Handler: resolveReportHandler,
	},
	"return_report": {
		MessageComponent: &discordgo.Button{
			Label:    "Вернуть",
			Style:    discordgo.PrimaryButton,
			CustomID: "return_report",
			Emoji: discordgo.ComponentEmoji{
				Name: "🔄",
			},
		},
		Handler: returnReportHandler,
	},
}

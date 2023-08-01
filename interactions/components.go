package interactions

import (
	"github.com/bwmarrin/discordgo"
)

type ComponentHandler func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)

var (
	ResolveReportButton = &discordgo.Button{
		CustomID: "resolve_report",
		Label:    "Рассмотрено",
		Style:    discordgo.SuccessButton,
		Emoji: discordgo.ComponentEmoji{
			Name: "✅",
		},
	}
	ReturnReportButton = &discordgo.Button{
		Label:    "Вернуть",
		Style:    discordgo.PrimaryButton,
		CustomID: "return_report",
		Emoji: discordgo.ComponentEmoji{
			Name: "🔄",
		},
	}
	ComponentHandlers = map[string]ComponentHandler{
		"resolve_report": resolveReportHandler,
		"return_report":  returnReportHandler,
		"remove_warning": removeWarningHandler,
	}
)

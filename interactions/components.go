package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/msg"
)

type ComponentHandler func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)

var (
	ResolveReportButton = &discordgo.Button{
		CustomID: "resolve_report",
		Label:    "Рассмотрено",
		Style:    discordgo.SecondaryButton,
		Emoji: discordgo.ComponentEmoji{
			Name: msg.CheckMarkEmoji.Name,
			ID:   msg.CheckMarkEmoji.ID,
		},
	}
	ComponentHandlers = map[string]ComponentHandler{
		"resolve_report": resolveReportHandler,
		"remove_warning": removeWarningHandler,
	}
)

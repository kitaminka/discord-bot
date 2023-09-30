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
		Emoji:    msg.ToComponentEmoji(msg.CheckMarkEmoji),
	}
	ComponentHandlers = map[string]ComponentHandler{
		"resolve_report": resolveReportHandler,
		"create_warning": createWarningHandler,
		"remove_warning": removeWarningHandler,
	}
)

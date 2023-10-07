package interactions

import (
	"fmt"
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
	ViewWarningsButton = func(userID string) *discordgo.Button {
		return &discordgo.Button{
			CustomID: fmt.Sprintf("view_warnings:%v", userID),
			Label:    "Предупреждения",
			Style:    discordgo.PrimaryButton,
		}
	}
	RemoveWarningsButton = func(userID string) *discordgo.Button {
		return &discordgo.Button{
			CustomID: fmt.Sprintf("remove_warnings:%v", userID),
			Label:    "Снять предупреждения",
			Style:    discordgo.DangerButton,
		}
	}
	ComponentHandlers = map[string]ComponentHandler{
		"resolve_report":  resolveReportHandler,
		"create_warning":  createWarningHandler,
		"remove_warning":  removeWarningHandler,
		"view_warnings":   viewWarningsButtonHandler,
		"remove_warnings": removeWarningsButtonHandler,
	}
)

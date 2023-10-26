package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/msg"
)

type ComponentHandler func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)

var (
	ResolveReportButton = func(senderID, userID, channelID, messageID string) *discordgo.Button {
		return &discordgo.Button{
			CustomID: fmt.Sprintf("resolve_report:%v:%v:%v:%v", senderID, userID, channelID, messageID),
			Label:    "Рассмотрено",
			Style:    discordgo.SuccessButton,
			Emoji:    msg.ToComponentEmoji(msg.CheckMarkEmoji),
		}
	}
	ReportWarningButton = func(senderID, userID, channelID, messageID string) *discordgo.Button {
		return &discordgo.Button{
			CustomID: fmt.Sprintf("report_warning:%v:%v:%v:%v", senderID, userID, channelID, messageID),
			Label:    "Выдать предупреждение",
			Style:    discordgo.SecondaryButton,
			Emoji:    msg.ToComponentEmoji(msg.ShieldCheckMarkEmoji),
		}
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
		"report_warning":  reportWarningButtonHandler,
	}
)

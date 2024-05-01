package logs

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/msg"
)

func LogUserMute(session *discordgo.Session, moderatorUser, targetUser *discordgo.User, reason string, muteUntil time.Time) {
	sendLogMessage(session, ModerationLog, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Мут выдан",
				Description: msg.StructuredText{
					Fields: []*msg.StructuredTextField{
						{
							Name:  "Окончание мута",
							Value: fmt.Sprintf("<t:%v:R>", muteUntil.Unix()),
						},
						{
							Name:  "Причина",
							Value: reason,
						},
						{
							Name:  "Пользователь",
							Value: msg.UserMention(targetUser),
						},
						{
							Name:  "Модератор",
							Value: msg.UserMention(moderatorUser),
						},
					},
				}.ToString(),
				Color: msg.DefaultEmbedColor,
			},
		},
	})
}

func LogUserUnmute(session *discordgo.Session, moderatorUser, targetUser *discordgo.User, muteUntil time.Time) {
	sendLogMessage(session, ModerationLog, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Мут снят",
				Description: msg.StructuredText{
					Fields: []*msg.StructuredTextField{
						{
							Name:  "Окончание мута",
							Value: fmt.Sprintf("<t:%v:R>", muteUntil.Unix()),
						},
						{
							Name:  "Пользователь",
							Value: msg.UserMention(targetUser),
						},
						{
							Name:  "Модератор",
							Value: msg.UserMention(moderatorUser),
						},
					},
				}.ToString(),
				Color: msg.DefaultEmbedColor,
			},
		},
	})
}

package logs

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/msg"
	"time"
)

func LogWarningCreation(session *discordgo.Session, moderatorUser, targetUser *discordgo.User, reason string, warnTime time.Time) {
	sendLogMessage(session, ModerationLog, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Предупреждение выдано",
				Description: msg.StructuredText{
					Fields: []*msg.StructuredTextField{
						{
							Name:  "Время выдачи",
							Value: fmt.Sprintf("<t:%v>", warnTime.Unix()),
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

func LogWarningRemoving(session *discordgo.Session, moderatorUser, targetUser *discordgo.User, reason string, warnTime time.Time) {
	sendLogMessage(session, ModerationLog, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Предупреждение снято",
				Description: msg.StructuredText{
					Fields: []*msg.StructuredTextField{
						{
							Name:  "Время выдачи",
							Value: fmt.Sprintf("<t:%v>", warnTime.Unix()),
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

func LogWarningResetting(session *discordgo.Session, moderatorUser, targetUser *discordgo.User) {
	sendLogMessage(session, ModerationLog, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Предупреждения сброшены",
				Description: msg.StructuredText{
					Fields: []*msg.StructuredTextField{
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

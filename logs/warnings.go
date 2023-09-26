package logs

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/msg"
	"log"
	"time"
)

func LogWarningCreation(session *discordgo.Session, moderatorUser, targetUser *discordgo.User, reason string, warnTime time.Time) {
	guild, err := db.GetGuild()
	if err != nil {
		log.Printf("Error getting guild: %v", err)
		return
	}

	_, err = session.ChannelMessageSendComplex(guild.ModerationLogChannelID, &discordgo.MessageSend{
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
	if err != nil {
		log.Printf("Error logging warning creation: %v", err)
		return
	}
}

func LogWarningRemoving(session *discordgo.Session, moderatorUser, targetUser *discordgo.User, reason string, warnTime time.Time) {
	guild, err := db.GetGuild()
	if err != nil {
		log.Printf("Error getting guild: %v", err)
		return
	}

	_, err = session.ChannelMessageSendComplex(guild.ModerationLogChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Предупреждение убрано",
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
	if err != nil {
		log.Printf("Error logging warning creation: %v", err)
		return
	}
}

func LogWarningResetting(session *discordgo.Session, moderatorUser, targetUser *discordgo.User) {
	guild, err := db.GetGuild()
	if err != nil {
		log.Printf("Error getting guild: %v", err)
		return
	}

	_, err = session.ChannelMessageSendComplex(guild.ModerationLogChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Предупреждение сброшены",
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
	if err != nil {
		log.Printf("Error logging warning creation: %v", err)
		return
	}
}

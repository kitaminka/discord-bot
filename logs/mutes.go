package logs

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/msg"
	"log"
	"time"
)

func LogUserMute(session *discordgo.Session, moderatorUser, targetUser *discordgo.User, reason string, muteUntil time.Time) {
	guild, err := db.GetGuild()
	if err != nil {
		log.Printf("Error getting guild: %v", err)
		return
	}

	_, err = session.ChannelMessageSendComplex(guild.ModerationLogChannelID, &discordgo.MessageSend{
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
	if err != nil {
		log.Printf("Error logging warning creation: %v", err)
		return
	}
}

func LogUserUnmute(session *discordgo.Session, moderatorUser, targetUser *discordgo.User, muteUntil time.Time) {
	guild, err := db.GetGuild()
	if err != nil {
		log.Printf("Error getting guild: %v", err)
		return
	}

	_, err = session.ChannelMessageSendComplex(guild.ModerationLogChannelID, &discordgo.MessageSend{
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
	if err != nil {
		log.Printf("Error logging warning creation: %v", err)
		return
	}
}

package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/msg"
	"log"
	"time"
)

func notifyUserWarning(session *discordgo.Session, userID string, warningTime time.Time, created bool, description string) {
	var embed *discordgo.MessageEmbed
	if created {
		embed = &discordgo.MessageEmbed{
			Title: "Вам выдано предупреждение",
			Description: msg.StructuredText{
				Text: description,
				Fields: []*msg.StructuredTextField{
					{
						Name:  "Время выдачи",
						Value: fmt.Sprintf("<t:%v>", warningTime.Unix()),
					},
					{
						Name:  "Время окончания",
						Value: fmt.Sprintf("<t:%v:R>", warningTime.Add(db.WarningDuration).Unix()),
					},
				},
			}.ToString(),
			Color: msg.DefaultEmbedColor,
		}
	} else {
		embed = &discordgo.MessageEmbed{
			Title: "С вас снято предупреждение",
			Description: msg.StructuredText{
				Fields: []*msg.StructuredTextField{
					{
						Name:  "Время выдачи",
						Value: fmt.Sprintf("<t:%v>", warningTime.Unix()),
					},
				},
			}.ToString(),
			Color: msg.DefaultEmbedColor,
		}
	}

	sendUserNotification(session, userID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	})
}

func notifyUserWarningReset(session *discordgo.Session, userID string) {
	sendUserNotification(session, userID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Ваши предупреждения сброшены",
				Description: "Все ваши предупреждения удалены.",
				Color:       msg.DefaultEmbedColor,
			},
		},
	})
}

func notifyUserMute(session *discordgo.Session, userID string, muteUntil time.Time, created bool, description string) {
	var embed *discordgo.MessageEmbed
	if created {
		embed = &discordgo.MessageEmbed{
			Title: "Вам выдан мут",
			Description: msg.StructuredText{
				Text: description,
				Fields: []*msg.StructuredTextField{
					{
						Name:  "Время окончания",
						Value: fmt.Sprintf("<t:%v:R>", muteUntil.Unix()),
					},
				},
			}.ToString(),
			Color: msg.DefaultEmbedColor,
		}
	} else {
		embed = &discordgo.MessageEmbed{
			Title: "С вас снят мут",
			Description: msg.StructuredText{
				Fields: []*msg.StructuredTextField{
					{
						Name:  "Время окончания",
						Value: fmt.Sprintf("<t:%v:R>", muteUntil.Unix()),
					},
				},
			}.ToString(),
			Color: msg.DefaultEmbedColor,
		}
	}

	sendUserNotification(session, userID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	})
}

func sendUserNotification(session *discordgo.Session, userID string, message *discordgo.MessageSend) {
	channel, err := session.UserChannelCreate(userID)
	if err != nil {
		log.Printf("Error creating user channel: %v", err)
		return
	}

	_, err = session.ChannelMessageSendComplex(channel.ID, message)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}
}

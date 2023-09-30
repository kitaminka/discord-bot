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
	channel, err := session.UserChannelCreate(userID)
	if err != nil {
		log.Printf("Error creating user channel: %v", err)
		return
	}

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

	_, err = session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	})
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}
}

func notifyUserWarningReset(session *discordgo.Session, userID string) {
	channel, err := session.UserChannelCreate(userID)
	if err != nil {
		log.Printf("Error creating user channel: %v", err)
		return
	}

	_, err = session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Ваши предупреждения сброшены",
				Description: "Все ваши предупреждения удалены.",
				Color:       msg.DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}
}

func notifyUserMute(session *discordgo.Session, userID string, muteUntil time.Time, created bool, description string) {
	channel, err := session.UserChannelCreate(userID)
	if err != nil {
		log.Printf("Error creating user channel: %v", err)
		return
	}

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

	_, err = session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	})
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}
}

package logs

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/msg"
)

func LogReputationChange(session *discordgo.Session, user, targetUser *discordgo.User, change int) {
	var description string
	switch change {
	case 1:
		description = fmt.Sprintf("%v поставил **лайк** %v", msg.UserMention(user), msg.UserMention(targetUser))
	case -1:
		description = fmt.Sprintf("%v поставил **дизлайк** %v", msg.UserMention(user), msg.UserMention(targetUser))
	default:
		description = fmt.Sprintf("%v изменил репутцию пользователя %v на **%v**.", msg.UserMention(user), msg.UserMention(targetUser), change)
	}

	sendLogMessage(session, ReputationLog, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Изменение репутации",
				Description: description,
				Color:       msg.DefaultEmbedColor,
			},
		},
	})
}

func LogReputationSetting(session *discordgo.Session, moderatorUser, targetUser *discordgo.User, value int) {
	sendLogMessage(session, ReputationLog, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Установка репутации",
				Description: fmt.Sprintf("%v установил репутцию пользователя %v на **%v**.", msg.UserMention(moderatorUser), msg.UserMention(targetUser), value),
				Color:       msg.DefaultEmbedColor,
			},
		},
	})
}

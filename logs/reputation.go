package logs

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/msg"
)

func LogReputationChange(session *discordgo.Session, guildID string, user, targetUser *discordgo.User, change int) error {
	guild, err := db.GetGuild()
	if err != nil {
		return err
	}

	var description string
	switch change {
	case 1:
		description = fmt.Sprintf("%v поставил лайк %v", msg.UserMention(user), msg.UserMention(targetUser))
	case -1:
		description = fmt.Sprintf("%v поставил дизлайк %v", msg.UserMention(user), msg.UserMention(targetUser))
	default:
		description = fmt.Sprintf("%v изменил репутцию пользователя %v на %v", msg.UserMention(user), msg.UserMention(targetUser), change)
	}

	_, err = session.ChannelMessageSendComplex(guild.ReputationLogChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Изменение репутации",
				Description: description,
				Color:       msg.DefaultEmbedColor,
			},
		},
	})
	return err
}

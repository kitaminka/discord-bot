package logs

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"log"
)

type logType int

const (
	ReputationLog logType = iota
	ModerationLog
)

func sendLogMessage(session *discordgo.Session, logType logType, message *discordgo.MessageSend) {
	guild, err := db.GetGuild()
	if err != nil {
		log.Printf("Error getting guild: %v", err)
		return
	}

	var channelID string
	switch logType {
	case ReputationLog:
		channelID = guild.ReputationLogChannelID
	case ModerationLog:
		channelID = guild.ModerationLogChannelID
	}

	_, err = session.ChannelMessageSendComplex(channelID, message)
	if err != nil {
		log.Printf("Error logging warning creation: %v", err)
		return
	}
}

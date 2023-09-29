package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/logs"
	"github.com/kitaminka/discord-bot/msg"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func getUserNextMuteDuration(user db.User) time.Duration {
	now := time.Now()

	if now.After(user.LastMuteTime.Add(ExtendedMutePeriod + time.Duration(int(MuteDuration)*user.MuteCount))) {
		return MuteDuration
	}

	return time.Duration((user.MuteCount + 1) * int(MuteDuration))
}

func muteChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Member.Permissions&discordgo.PermissionModerateMembers == 0 {
		InteractionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
		return
	}

	var (
		discordUser    *discordgo.User
		durationString string
		reasonString   string
	)

	for _, option := range interactionCreate.ApplicationCommandData().Options {
		switch option.Name {
		case "пользователь":
			discordUser = option.UserValue(session)
		case "длительность":
			durationString = option.StringValue()
		case "причина":
			reasonString = option.StringValue()
		}
	}

	if discordUser.Bot {
		InteractionRespondError(session, interactionCreate.Interaction, "Вы не можете выдать мут боту.")
		return
	}
	isModerator, err := isUserModerator(session, interactionCreate.Interaction, discordUser)
	if err != nil {
		InteractionRespondError(session, interactionCreate.Interaction, "Произошла ошибка при выдаче муиа. Свяжитесь с администрацией.")
		log.Printf("Error checking if user is moderator: %v", err)
		return
	}
	if isModerator {
		InteractionRespondError(session, interactionCreate.Interaction, "Вы не можете выдать мут себе или другому модератору.")
		return
	}

	durationString = strings.ReplaceAll(durationString, " ", "")
	duration, err := time.ParseDuration(durationString)
	if err != nil {
		InteractionRespondError(session, interactionCreate.Interaction, "Неверный формат длительности.")
		return
	}
	if duration < 0 {
		InteractionRespondError(session, interactionCreate.Interaction, "Длительность мута не может быть отрицательной.")
		return
	}

	until := time.Now().Add(duration)

	reasonIndex, err := strconv.Atoi(reasonString)
	if err != nil {
		InteractionRespondError(session, interactionCreate.Interaction, "Произошла ошибка при выдаче мута. Свяжитесь с администрацией.")
		log.Printf("Error parsing reason index: %v", err)
		return
	}

	reason := Reasons[reasonIndex]

	err = session.GuildMemberTimeout(interactionCreate.GuildID, discordUser.ID, &until, discordgo.WithAuditLogReason(url.QueryEscape(reason.Name)))
	if err != nil {
		InteractionRespondError(session, interactionCreate.Interaction, "Произошла ошибка при выдаче мута. Свяжитесь с администрацией.")
		log.Printf("Error muting member: %v", err)
		return
	}

	err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Мут выдан",
					Description: msg.StructuredText{
						Text: "Мут успешно выдан.",
						Fields: []*msg.StructuredTextField{
							{
								Name:  "Время окончания",
								Value: fmt.Sprintf("<t:%v:R>", until.Unix()),
							},
							{
								Name:  "Причина",
								Value: reason.Name,
							},
							{
								Name:  "Пользователь",
								Value: msg.UserMention(discordUser),
							},
							{
								Name:  "Модератор",
								Value: msg.UserMention(interactionCreate.Member.User),
							},
						},
					}.ToString(),
					Color: msg.DefaultEmbedColor,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}
	logs.LogUserMute(session, interactionCreate.Member.User, discordUser, reason.Name, until)
}

func unmuteChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Member.Permissions&discordgo.PermissionModerateMembers == 0 {
		InteractionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
		return
	}

	var discordUser *discordgo.User

	for _, option := range interactionCreate.ApplicationCommandData().Options {
		switch option.Name {
		case "пользователь":
			discordUser = option.UserValue(session)
		}
	}

	if discordUser.Bot {
		InteractionRespondError(session, interactionCreate.Interaction, "Вы не можете снять мут с бота.")
		return
	}

	member, err := session.GuildMember(interactionCreate.GuildID, discordUser.ID)
	if err != nil {
		InteractionRespondError(session, interactionCreate.Interaction, "Произошла ошибка при снятии мута. Свяжитесь с администрацией.")
		log.Printf("Error getting member: %v", err)
		return
	}

	muteUntil := member.CommunicationDisabledUntil

	if muteUntil == nil || time.Now().After(*muteUntil) {
		InteractionRespondError(session, interactionCreate.Interaction, "У пользователя нет мута.")
		return
	}

	err = db.ResetUserMuteCount(discordUser.ID)
	if err != nil {
		InteractionRespondError(session, interactionCreate.Interaction, "Произошла ошибка при снятии мута. Свяжитесь с администрацией.")
		log.Printf("Error resetting user mute count: %v", err)
		return
	}

	err = session.GuildMemberTimeout(interactionCreate.GuildID, discordUser.ID, nil, discordgo.WithAuditLogReason(url.QueryEscape(fmt.Sprintf("Снятие мута от %v", interactionCreate.Member.User.Username))))
	if err != nil {
		InteractionRespondError(session, interactionCreate.Interaction, "Произошла ошибка при снятии мута. Свяжитесь с администрацией.")
		log.Printf("Error unmuting member: %v", err)
		return
	}

	err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Мут снят",
					Description: msg.StructuredText{
						Text: "Мут успешно снят.",
						Fields: []*msg.StructuredTextField{
							{
								Name:  "Пользователь",
								Value: msg.UserMention(discordUser),
							},
							{
								Name:  "Время окончания",
								Value: fmt.Sprintf("<t:%v:R>", muteUntil.Unix()),
							},
						},
					}.ToString(),
					Color: msg.DefaultEmbedColor,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	logs.LogUserUnmute(session, interactionCreate.Member.User, discordUser, *muteUntil)
}

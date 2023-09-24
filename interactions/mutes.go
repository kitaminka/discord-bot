package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/msg"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func getUserNextMuteDuration(user db.User, warnings []db.Warning) time.Duration {
	if len(warnings) < MuteWarningsCount {
		return 0
	}

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
		return
	}

	reason := Reasons[reasonIndex]

	err = session.GuildMemberTimeout(interactionCreate.GuildID, discordUser.ID, &until, discordgo.WithAuditLogReason(url.QueryEscape(reason.Name)))
	if err != nil {
		InteractionRespondError(session, interactionCreate.Interaction, "Произошла ошибка при выдаче мута. Свяжитесь с администрацией.")
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
								Name:  "Мут истекает",
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
}

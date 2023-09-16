package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/msg"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strconv"
	"time"
)

func warnChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Member.Permissions&discordgo.PermissionModerateMembers == 0 {
		InteractionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
		return
	}

	var (
		discordUser  *discordgo.User
		reasonString string
	)

	for _, option := range interactionCreate.ApplicationCommandData().Options {
		switch option.Name {
		case "пользователь":
			discordUser = option.UserValue(session)
		case "причина":
			reasonString = option.StringValue()
		}
	}

	if discordUser.Bot {
		InteractionRespondError(session, interactionCreate.Interaction, "Вы не можете выдать предупреждение боту.")
		return
	}

	err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}

	warnTime := time.Now()

	reasonIndex, err := strconv.Atoi(reasonString)
	if err != nil {
		InteractionRespondError(session, interactionCreate.Interaction, "Произошла ошибка при выдаче мута. Свяжитесь с администрацией.")
		return
	}

	reason := msg.Reasons[reasonIndex]

	err = db.CreateWarning(db.Warning{
		Time:        warnTime,
		Reason:      reason.Name,
		UserID:      discordUser.ID,
		ModeratorID: interactionCreate.Member.User.ID,
	})

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title: "Предупреждение выдано",
				Description: msg.StructuredText{
					Text: "Предупреждение успешно выдано.",
					Fields: []*msg.StructuredTextField{
						{
							Name:  "Время выдачи",
							Value: fmt.Sprintf("<t:%v>", warnTime.Unix()),
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
	})
}
func warnMessageCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Member.Permissions&discordgo.PermissionModerateMembers == 0 {
		InteractionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
		return
	}

	message := interactionCreate.ApplicationCommandData().Resolved.Messages[interactionCreate.ApplicationCommandData().TargetID]

	if message.Author.Bot {
		InteractionRespondError(session, interactionCreate.Interaction, "Вы не можете выдать предупреждение на сообщение бота.")
		return
	}

	perms, err := session.UserChannelPermissions(message.Author.ID, interactionCreate.ChannelID)
	if err != nil {
		InteractionRespondError(session, interactionCreate.Interaction, "Произошла ошибка при выдаче предупреждения. Свяжитесь с администрацией.")
		log.Printf("Error getting user permissions: %v", err)
		return
	}

	guild, err := session.Guild(interactionCreate.GuildID)
	if err != nil {
		InteractionRespondError(session, interactionCreate.Interaction, "Произошла ошибка при выдаче предупреждения. Свяжитесь с администрацией.")
		log.Printf("Error getting guild: %v", err)
		return
	}

	if perms&discordgo.PermissionModerateMembers != 0 && interactionCreate.Member.User.ID == guild.OwnerID {
		InteractionRespondError(session, interactionCreate.Interaction, "Вы не можете выдать предупреждение себе или другому модератору.")
		return
	}

	err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Выдача предупреждения",
					Description: msg.StructuredText{
						Text: "Выберите причину предупреждения.",
						Fields: []*msg.StructuredTextField{
							{
								Name:  "Пользователь",
								Value: msg.UserMention(message.Author),
							},
							{
								Name:  "Сообщение",
								Value: fmt.Sprintf("https://discord.com/channels/%v/%v/%v", interactionCreate.GuildID, interactionCreate.ChannelID, message.ID),
							},
						},
					}.ToString(),
					Color: msg.DefaultEmbedColor,
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						createWarnSelectMenu(),
					},
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
func createWarnSelectMenu() discordgo.SelectMenu {
	var selectMenuOptions []discordgo.SelectMenuOption

	for i, reason := range msg.Reasons {
		selectMenuOptions = append(selectMenuOptions, discordgo.SelectMenuOption{
			Label:       reason.Name,
			Value:       strconv.Itoa(i),
			Description: reason.Description,
			Emoji: discordgo.ComponentEmoji{
				Name: msg.ReportEmoji.Name,
				ID:   msg.ReportEmoji.ID,
			},
		})
	}

	return discordgo.SelectMenu{
		MenuType:    discordgo.StringSelectMenu,
		CustomID:    "warn",
		Placeholder: "Выберите причину",
		Options:     selectMenuOptions,
	}
}

func createRemWarnSelectMenu(session *discordgo.Session, warnings []db.Warning) (discordgo.SelectMenu, error) {
	var selectMenuOptions []discordgo.SelectMenuOption

	for i, warn := range warnings {
		moderatorDiscordUser, err := session.User(warn.ModeratorID)
		if err != nil {
			return discordgo.SelectMenu{}, err
		}

		selectMenuOptions = append(selectMenuOptions, discordgo.SelectMenuOption{
			Label:       fmt.Sprintf("Предупреждение #%v от %v", i+1, moderatorDiscordUser.Username),
			Value:       warn.ID.Hex(),
			Description: fmt.Sprintf("Причина: %v", warn.Reason),
			Emoji: discordgo.ComponentEmoji{
				Name: msg.ReportEmoji.Name,
				ID:   msg.ReportEmoji.ID,
			},
		})
	}

	return discordgo.SelectMenu{
		MenuType:    discordgo.StringSelectMenu,
		CustomID:    "remove_warning",
		Placeholder: "Выберите предупреждение",
		Options:     selectMenuOptions,
	}, nil
}

func remWarnChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
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
		InteractionRespondError(session, interactionCreate.Interaction, "Вы не можете снять предупреждение с бота.")
		return
	}

	err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}

	warnings, err := db.GetUserWarnings(discordUser.ID)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при снятии предупреждений. Свяжитесь с администрацией.")
		log.Printf("Error getting user: %v", err)
		return
	}

	if len(warnings) == 0 {
		_, _ = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Предупреждения отсутствуют",
					Description: "У данного пользователя нет предупреждений.",
					Color:       msg.DefaultEmbedColor,
				},
			},
		})
		return
	}

	remWarnSelectMenu, err := createRemWarnSelectMenu(session, warnings)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при снятии предупреждений. Свяжитесь с администрацией.")
		log.Printf("Error creating select menu: %v", err)
		return
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "Снятие предупреждений",
				Description: "Выберите предупреждение, которое вы хотите снять.",
				Color:       msg.DefaultEmbedColor,
			},
		},
		Components: &[]discordgo.MessageComponent{
			&discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					remWarnSelectMenu,
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error editing interaction response: %v", err)
		return
	}
}

func removeWarningHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Member.Permissions&discordgo.PermissionModerateMembers == 0 {
		InteractionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
		return
	}

	componentValue := interactionCreate.MessageComponentData().Values[0]

	warnID, err := primitive.ObjectIDFromHex(componentValue)
	if err != nil {
		InteractionRespondError(session, interactionCreate.Interaction, "Произошла ошибка при снятии предупреждения. Свяжитесь с администрацией.")
		log.Printf("Error creating object ID: %v", err)
		return
	}

	err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}

	warning, err := db.RemoveWarning(warnID)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при снятии предупреждения. Свяжитесь с администрацией.")
		log.Printf("Error removing warning: %v", err)
		return
	}
	discordUser, err := session.User(warning.UserID)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при снятии предупреждения. Свяжитесь с администрацией.")
		log.Printf("Error getting user: %v", err)
		return
	}
	moderatorDiscordUser, err := session.User(warning.ModeratorID)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при снятии предупреждения. Свяжитесь с администрацией.")
		log.Printf("Error getting user: %v", err)
		return
	}
	warnings, err := db.GetUserWarnings(warning.UserID)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при снятии предупреждения. Свяжитесь с администрацией.")
		log.Printf("Error getting user warnings: %v", err)
		return
	}

	if len(warnings) == 0 {
		_, _ = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Предупреждения отсутствуют",
					Description: "У данного пользователя нет предупреждений.",
					Color:       msg.DefaultEmbedColor,
				},
			},
			Components: &[]discordgo.MessageComponent{},
		})
	} else {
		remWarnSelectMenu, err := createRemWarnSelectMenu(session, warnings)
		if err != nil {
			log.Printf("Error creating select menu: %v", err)
		} else {
			_, _ = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
				Embeds: &interactionCreate.Message.Embeds,
				Components: &[]discordgo.MessageComponent{
					&discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							remWarnSelectMenu,
						},
					},
				},
			})
		}
	}

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, false, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Предупреждение снято",
				Description: msg.StructuredText{
					Text: "Предупреждение снято.",
					Fields: []*msg.StructuredTextField{
						{
							Name:  "ID",
							Value: fmt.Sprintf("`%v`", warning.ID.Hex()),
						},
						{
							Name:  "Время выдачи",
							Value: fmt.Sprintf("<t:%v>", warning.Time.Unix()),
						},
						{
							Name:  "Причина",
							Value: warning.Reason,
						},
						{
							Name:  "Пользователь",
							Value: msg.UserMention(discordUser),
						},
						{
							Name:  "Модератор",
							Value: msg.UserMention(moderatorDiscordUser),
						},
					},
				}.ToString(),
				Color: msg.DefaultEmbedColor,
			},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		log.Printf("Error creating followup message: %v", err)
		return
	}
}

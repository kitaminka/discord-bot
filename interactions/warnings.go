package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/msg"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func warnChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Member.Permissions&discordgo.PermissionModerateMembers == 0 {
		interactionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
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
		interactionRespondError(session, interactionCreate.Interaction, "Вы не можете выдать предупреждение боту.")
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

	err = db.AddUserWarning(db.Warning{
		Time:        warnTime,
		UserID:      discordUser.ID,
		ModeratorID: interactionCreate.Member.User.ID,
	})

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title: "Предупреждение выдано",
				Description: msg.StructuredDescription{
					Text: "Предупреждение успешно выдано.",
					Fields: []*msg.StructuredDescriptionField{
						{
							Name:  "Время выдачи",
							Value: fmt.Sprintf("<t:%v>", warnTime.Unix()),
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
			Description: fmt.Sprintf("Айди: %v", warn.ID.Hex()),
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
		interactionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
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
		interactionRespondError(session, interactionCreate.Interaction, "Вы не можете снять предупреждение с бота.")
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
		interactionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
		return
	}

	componentValue := interactionCreate.MessageComponentData().Values[0]

	warnID, err := primitive.ObjectIDFromHex(componentValue)
	if err != nil {
		interactionRespondError(session, interactionCreate.Interaction, "Произошла ошибка при снятии предупреждения. Свяжитесь с администрацией.")
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
				Description: msg.StructuredDescription{
					Text: "Предупреждение снято.",
					Fields: []*msg.StructuredDescriptionField{
						{
							Name:  "ID",
							Value: warning.ID.Hex(),
						},
						{
							Name:  "Время выдачи",
							Value: fmt.Sprintf("<t:%v>", warning.Time.Unix()),
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

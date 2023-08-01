package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/msg"
	"log"
	"strings"
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

	err = db.AddUserWarn(discordUser.ID, db.Warning{
		Time:        time.Now(),
		ModeratorID: interactionCreate.Member.User.ID,
	})

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title: "Предупреждение выдано",
				Description: msg.StructuredDescription{
					Text: "Предупреждение успешно выдано.",
				}.ToString(),
				Color: msg.DefaultEmbedColor,
			},
		},
	})
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

	user, err := db.GetUser(discordUser.ID)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при снятии предупреждения. Свяжитесь с администрацией.")
		log.Printf("Error getting user: %v", err)
		return
	}

	if len(user.Warnings) == 0 {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "У данного пользователя нет предупреждений.")
		return
	}

	var selectMenuOptions []discordgo.SelectMenuOption

	for i, warn := range user.Warnings {
		selectMenuOptions = append(selectMenuOptions, discordgo.SelectMenuOption{
			Label:       fmt.Sprintf("Предупреждение #%v", i+1),
			Value:       fmt.Sprintf("%v_%v", user.ID, warn.Time),
			Description: "Пред",
			Emoji: discordgo.ComponentEmoji{
				Name: "report",
				ID:   "1128761333014347890",
			},
		})
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
					discordgo.SelectMenu{
						MenuType:    discordgo.StringSelectMenu,
						CustomID:    "remove_warning",
						Placeholder: "Выберите предупреждение",
						Options:     selectMenuOptions,
					},
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

	userID := strings.Split(componentValue, "_")[0]
	warnTimeString := strings.Split(componentValue, "_")[1]

	warnTime, _ := time.Parse(time.RFC3339, warnTimeString)

	fmt.Println(warnTime)

	err := db.RemoveUserWarn(userID, db.Warning{
		Time: warnTime,
	})

	fmt.Println(err)
}

package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/logs"
	"github.com/kitaminka/discord-bot/msg"
	"log"
	"strconv"
	"time"
)

func likeUserCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	reputationCommandHandler(session, interactionCreate, true)
}
func dislikeUserCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	reputationCommandHandler(session, interactionCreate, false)
}

func likeChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	reputationCommandHandler(session, interactionCreate, true)
}
func dislikeChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	reputationCommandHandler(session, interactionCreate, false)
}

// Used for user command and chat command
func reputationCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate, like bool) {
	var action string
	var reputationChange int
	var discordUser *discordgo.User

	if like {
		action = "лайк"
		reputationChange = 1
	} else {
		action = "дизлайк"
		reputationChange = -1
	}

	if len(interactionCreate.ApplicationCommandData().Options) == 0 {
		var err error
		discordUser, err = session.User(interactionCreate.ApplicationCommandData().TargetID)
		if err != nil {
			interactionRespondError(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
			log.Printf("Error getting user: %v", err)
			return
		}
	} else {
		discordUser = interactionCreate.ApplicationCommandData().Options[0].UserValue(session)
	}

	if interactionCreate.Member.User.ID == discordUser.ID {
		interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Вы не можете поставить %v самому себе.", action))
		return
	}

	if discordUser.Bot {
		interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Вы не можете поставить %v боту.", action))
		return
	}

	_, err := session.GuildMember(interactionCreate.GuildID, discordUser.ID)
	if err != nil {
		interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Вы не можете поставить %v пользователю, который не находится на сервере.", action))
		return
	}

	err = session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}

	user, err := db.GetUser(interactionCreate.Member.User.ID)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
		log.Printf("Error getting reputation delay: %v", err)
		return
	}
	if !time.Now().After(user.ReputationDelay) {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, fmt.Sprintf("Вы сможете сделать это <t:%v:R>.", user.ReputationDelay.Unix()))
		return
	}

	err = db.UpdateUserReputationDelay(user.ID)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
		log.Printf("Error updating reputation delay: %v", err)
		return
	}

	err = db.ChangeUserReputation(discordUser.ID, reputationChange)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
		log.Printf("Error changing user reputation: %v", err)
		err = db.ResetUserReputationDelay(user.ID)
		if err != nil {
			log.Printf("Error clearing reputation delay: %v", err)
		}
		return
	}

	err = logs.LogReputationChange(session, interactionCreate.Member.User, discordUser, reputationChange)
	if err != nil {
		log.Printf("Error logging reputation change: %v", err)
	}

	var title string
	if like {
		title = fmt.Sprintf("%v Лайк", msg.LikeEmoji)
	} else {
		title = fmt.Sprintf("%v Дизлайк", msg.DislikeEmoji)
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       title,
				Description: fmt.Sprintf("Вы поставили %v пользователю %v.", action, msg.UserMention(discordUser)),
				Color:       msg.DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Printf("Error editing interaction response: %v", err)
		return
	}
}

func topReputationChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}

	users, err := db.GetUserReputationTop()
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
		log.Printf("Error getting reputation top: %v", err)
		return
	}

	var fields []*discordgo.MessageEmbedField
	for i, user := range *users {
		var discordUser *discordgo.User
		discordUser, err = session.User(user.ID)
		if err != nil {
			interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
			log.Printf("Error getting user: %v", err)
			return
		}

		var PlaceEmoji msg.Emoji

		switch i {
		case 0:
			PlaceEmoji = msg.FirstEmoji
		case 1:
			PlaceEmoji = msg.SecondEmoji
		case 2:
			PlaceEmoji = msg.ThirdEmoji
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%v #%v. %v", PlaceEmoji, i+1, discordUser.Username),
			Value: fmt.Sprintf("%v **Репутация**: %v", msg.ReputationEmoji, user.Reputation),
		})
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:  "Топ пользователей по репутации",
				Fields: fields,
				Color:  msg.DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Printf("Error editing interaction response: %v", err)
		return
	}
}

func setReputationChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		interactionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
		return
	}

	var discordUser *discordgo.User
	var reputation int

	for _, option := range interactionCreate.ApplicationCommandData().Options {
		switch option.Name {
		case "пользователь":
			discordUser = option.UserValue(session)
		case "репутация":
			reputation = int(option.IntValue())
		}
	}

	if discordUser == nil {
		interactionRespondError(session, interactionCreate.Interaction, "Не указан пользователь.")
		return
	} else if reputation == 0 {
		interactionRespondError(session, interactionCreate.Interaction, "Не указана репутация.")
		return
	}

	if discordUser.Bot {
		interactionRespondError(session, interactionCreate.Interaction, "Вы не можете изменить репутацию бота.")
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

	err = db.SetUserReputation(discordUser.ID, reputation)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при изменении репутации пользователя. Свяжитесь с администрацией.")
		log.Printf("Error setting user reputation: %v", err)
		return
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title: "Репутация изменена",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Пользователь",
						Value: msg.UserMention(discordUser),
					},
					{
						Name:  "Репутация",
						Value: strconv.Itoa(reputation),
					},
				},
				Color: msg.DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Printf("Error editing interaction response: %v", err)
		return
	}
}

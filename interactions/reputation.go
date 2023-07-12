package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/logs"
	"github.com/kitaminka/discord-bot/msg"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
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
	var targetUser *discordgo.User

	if like {
		action = "лайк"
		reputationChange = 1
	} else {
		action = "дизлайк"
		reputationChange = -1
	}

	if len(interactionCreate.ApplicationCommandData().Options) == 0 {
		var err error
		targetUser, err = session.User(interactionCreate.ApplicationCommandData().TargetID)
		if err != nil {
			interactionRespondError(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
			log.Printf("Error getting user: %v", err)
			return
		}
	} else {
		targetUser = interactionCreate.ApplicationCommandData().Options[0].UserValue(session)
	}

	if interactionCreate.Member.User.ID == targetUser.ID {
		interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Вы не можете поставить %v самому себе.", action))
		return
	}

	if targetUser.Bot {
		interactionRespondError(session, interactionCreate.Interaction, fmt.Sprintf("Вы не можете поставить %v боту.", action))
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

	user, err := db.GetUser(interactionCreate.Member.User.ID)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
		log.Printf("Error getting reputation delay: %v", err)
		return
	}
	if !time.Now().After(user.ReputationDelay) {
		followupErrorMessageCreate(session, interactionCreate.Interaction, fmt.Sprintf("Вы сможете сделать это <t:%v:R>.", user.ReputationDelay.Unix()))
		return
	}

	err = db.UpdateUserReputationDelay(user.ID)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
		log.Printf("Error updating reputation delay: %v", err)
		return
	}

	err = db.ChangeUserReputation(targetUser.ID, reputationChange)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
		log.Printf("Error changing user reputation: %v", err)
		err = db.ResetUserReputationDelay(user.ID)
		if err != nil {
			log.Printf("Error clearing reputation delay: %v", err)
		}
		return
	}

	err = logs.LogReputationChange(session, interactionCreate.GuildID, interactionCreate.Member.User, targetUser, reputationChange)
	if err != nil {
		log.Printf("Error logging reputation change: %v", err)
	}

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       cases.Title(language.Russian).String(action),
				Description: fmt.Sprintf("Вы поставили %v пользователю %v.", action, msg.UserMention(targetUser)),
				Color:       msg.DefaultEmbedColor,
			},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		log.Printf("Error creating followup message: %v", err)
		return
	}
}

func topReputationChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
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

	users, err := db.GetUserReputationTop()
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
		log.Printf("Error getting reputation top: %v", err)
		return
	}

	var fields []*discordgo.MessageEmbedField
	for i, user := range *users {
		var discordUser *discordgo.User
		discordUser, err = session.User(user.ID)
		if err != nil {
			followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
			log.Printf("Error getting user: %v", err)
			return
		}

		var PlaceEmoji string

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

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:  "Топ пользователей по репутации",
				Fields: fields,
				Color:  msg.DefaultEmbedColor,
			},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		log.Printf("Error creating followup message: %v", err)
		return
	}
}

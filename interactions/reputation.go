package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
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

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       cases.Title(language.Russian).String(action),
				Description: fmt.Sprintf("Вы поставили %v пользователю %v.", action, targetUser.Mention()),
				Color:       DefaultEmbedColor,
			},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		log.Printf("Error creating followup message: %v", err)
		return
	}
}

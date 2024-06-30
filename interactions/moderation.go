package interactions

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/msg"
)

const (
	MuteDuration       = 36 * time.Hour
	ExtendedMutePeriod = 48 * time.Hour
	MuteWarningsCount  = 3
)

func isUserModerator(session *discordgo.Session, interaction *discordgo.Interaction, discordUser *discordgo.User) (bool, error) {
	perms, err := session.UserChannelPermissions(discordUser.ID, interaction.ChannelID)
	if err != nil {
		return false, err
	}

	guild, err := session.Guild(interaction.GuildID)
	if err != nil {
		return false, err
	}

	return perms&discordgo.PermissionModerateMembers != 0 || discordUser.ID == guild.OwnerID, nil
}

func banChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Member.Permissions&discordgo.PermissionModerateMembers == 0 {
		InteractionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
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

	guild, err := db.GetGuild()
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
		log.Printf("Error getting guild: %v", err)
		return
	}

	supremeModerator := false

	for _, role := range interactionCreate.Member.Roles {
		if role == guild.SupremeModeratorRoleID {
			supremeModerator = true
		}
	}

	if !supremeModerator {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
		return
	}

	moderatorUser, err := db.GetUser(interactionCreate.Member.User.ID)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
		log.Printf("Error getting guild: %v", err)
		return
	}

	discordUser := interactionCreate.ApplicationCommandData().Options[0].UserValue(session)

	isModerator, err := isUserModerator(session, interactionCreate.Interaction, discordUser)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
		log.Printf("Error checking if user is moderator: %v", err)
		return
	}
	if isModerator {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Вы не можете забанить другого модератора.")
		return
	}

	if interactionCreate.Member.User.ID == discordUser.ID {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Вы не можете забанить самого себя.")
		return
	}

	if discordUser.Bot {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Вы не можете забанить бота.")
		return
	}

	_, err = session.GuildMember(interactionCreate.GuildID, discordUser.ID)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Вы не можете забанить человека, которого нет на сервере.")
		return
	}

	if !time.Now().After(moderatorUser.BanDelay) {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, fmt.Sprintf("Вы сможете сделать это <t:%v:R>.", moderatorUser.BanDelay.Unix()))
		return
	}

	err = db.UpdateUserBanDelay(moderatorUser.ID)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
		log.Printf("Error updating ban delay: %v", err)
		return
	}

	err = session.GuildBanCreateWithReason(interactionCreate.GuildID, discordUser.ID, fmt.Sprintf("Забанен модератором %v", interactionCreate.Member.User.Username), 7)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
		log.Printf("Error banning user: %v", err)
		return
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Description: fmt.Sprintf("%v был забанен", discordUser.Mention()),
				Color:       msg.DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка. Свяжитесь с администрацией.")
		log.Printf("Error editing interaction response: %v", err)
		return
	}
}

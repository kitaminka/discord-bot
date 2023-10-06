package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"github.com/kitaminka/discord-bot/msg"
	"log"
	"strconv"
	"time"
)

var (
	AdministratorPermission = int64(discordgo.PermissionAdministrator)
	ModeratorPermission     = int64(discordgo.PermissionModerateMembers)
)

// Used for user command and chat command
func profileCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
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

	var member *discordgo.Member

	if len(interactionCreate.ApplicationCommandData().TargetID) != 0 {
		member, err = session.GuildMember(interactionCreate.GuildID, interactionCreate.ApplicationCommandData().TargetID)
		if err != nil {
			interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при получении профиля пользователя. Свяжитесь с администрацией.")
			log.Printf("Error getting member: %v", err)
			return
		}
	} else if len(interactionCreate.ApplicationCommandData().Options) == 0 {
		member = interactionCreate.Member
	} else {
		discordUser := interactionCreate.ApplicationCommandData().Options[0].UserValue(session)

		member, err = session.GuildMember(interactionCreate.GuildID, discordUser.ID)
		if err != nil {
			interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при получении профиля пользователя. Свяжитесь с администрацией.")
			log.Printf("Error getting member: %v", err)
			return
		}
	}

	if member.User.Bot {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Вы не можете просмотреть профиль бота.")
		return
	}

	user, err := db.GetUser(member.User.ID)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при получении профиля пользователя. Свяжитесь с администрацией.")
		log.Printf("Error getting user: %v", err)
		return
	}

	var components []discordgo.MessageComponent

	embeds := []*discordgo.MessageEmbed{
		{
			Title: fmt.Sprintf("%v Профиль пользователя %v", msg.UserEmoji.MessageFormat(), member.User.Username),
			Description: msg.StructuredText{
				Fields: []*msg.StructuredTextField{
					{
						Emoji: msg.UsernameEmoji,
						Name:  "Пользователь",
						Value: member.Mention(),
					},
					{
						Emoji: msg.JoinEmoji,
						Name:  "Присоединился к серверу",
						Value: fmt.Sprintf("<t:%v:R>", member.JoinedAt.Unix()),
					},
					{
						Emoji: msg.ReputationEmoji,
						Name:  "Репутация",
						Value: strconv.Itoa(user.Reputation),
					},
					{
						Emoji: msg.ReportEmoji,
						Name:  "Отправлено репортов",
						Value: strconv.Itoa(user.ReportsSentCount),
					},
				},
			}.ToString(),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: member.AvatarURL(""),
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("ID: %v", member.User.ID),
			},
			Color: msg.DefaultEmbedColor,
		},
	}

	isModerator, err := isUserModerator(session, interactionCreate.Interaction, interactionCreate.Member.User)
	if err != nil {
		interactionResponseErrorEdit(session, interactionCreate.Interaction, "Произошла ошибка при получении профиля пользователя. Свяжитесь с администрацией.")
		log.Printf("Error checking if user is moderator: %v", err)
		return
	}

	if isModerator {
		var description msg.StructuredText

		if user.MuteCount > 0 && time.Now().Before(user.LastMuteTime.Add(ExtendedMutePeriod*time.Duration(user.MuteCount))) {
			description.Fields = append(description.Fields, &msg.StructuredTextField{
				Emoji: msg.TextChannelEmoji,
				Name:  "Количество предыдущих мутов",
				Value: strconv.Itoa(user.MuteCount),
			})
		}
		if time.Now().Before(user.LastMuteTime.Add(ExtendedMutePeriod * time.Duration(user.MuteCount))) {
			description.Fields = append(description.Fields, &msg.StructuredTextField{
				Emoji: msg.ShieldCheckMarkEmoji,
				Name:  "Окончание повышенного времени мута",
				Value: fmt.Sprintf("<t:%v:R>", user.LastMuteTime.Add((ExtendedMutePeriod+MuteDuration)*time.Duration(user.MuteCount)).Unix()),
			})
		}
		if member.CommunicationDisabledUntil != nil && time.Now().Before(*member.CommunicationDisabledUntil) {
			description.Fields = append(description.Fields, &msg.StructuredTextField{
				Emoji: msg.ShieldCheckMarkEmoji,
				Name:  "Время окончания мута",
				Value: fmt.Sprintf("<t:%v:R>", member.CommunicationDisabledUntil.Unix()),
			})
		}
		if len(description.Fields) == 0 {
			description.Text = "Нет информации."
		}

		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					ViewWarningsButton,
					RemoveWarningsButton,
				},
			},
		}

		embeds = append(embeds, &discordgo.MessageEmbed{
			Title:       "Модераторская информация",
			Description: description.ToString(),
			Color:       msg.DefaultEmbedColor,
		})
	}

	_, err = session.InteractionResponseEdit(interactionCreate.Interaction, &discordgo.WebhookEdit{
		Embeds:     &embeds,
		Components: &components,
	})
	if err != nil {
		log.Printf("Error editing interaction response: %v", err)
		return
	}
}

func topChatCommandHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	switch interactionCreate.ApplicationCommandData().Options[0].Name {
	case "reputation":
		topReputationChatCommandHandler(session, interactionCreate)
	default:
		InteractionRespondError(session, interactionCreate.Interaction, "Неизвестная подкоманда.")
	}
}

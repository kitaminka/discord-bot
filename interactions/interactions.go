package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"log"
)

const (
	DefaultEmbedColor = 14546431
	ErrorEmbedColor   = 16711680
)

var AdministratorPermission = int64(discordgo.PermissionAdministrator)

func UpdateGuildHandler(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	if interactionCreate.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		interactionRespondError(session, interactionCreate.Interaction, "Извините, но у вас нет прав на использование этой команды.")
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

	var fields []*discordgo.MessageEmbedField
	server := db.Guild{
		ID: interactionCreate.GuildID,
	}

	for _, option := range interactionCreate.ApplicationCommandData().Options {
		switch option.Name {
		case "report-channel":
			channel := option.ChannelValue(session)
			server.ReportChannelID = channel.ID
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  "Канал для репортов",
				Value: channel.Mention(),
			})
		case "resolved-report-channel":
			channel := option.ChannelValue(session)
			server.ResoledReportChannelID = channel.ID
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  "Канал для рассмотренных репортов",
				Value: channel.Mention(),
			})
		}
	}

	err = db.UpdateGuild(server)
	if err != nil {
		followupErrorMessageCreate(session, interactionCreate.Interaction, "Произошла ошибка при обновлении настреок сервера.")
		log.Printf("Error updating guild: %v", err)
		return
	}

	_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Настройки сервера обновлены",
				Description: "Настройки сервера были успешно обновлены.",
				Fields:      fields,
				Color:       DefaultEmbedColor,
			},
		},
	})
	if err != nil {
		log.Printf("Error creating followup message: %v", err)
		return
	}
}

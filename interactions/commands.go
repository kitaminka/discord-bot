package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kitaminka/discord-bot/db"
	"log"
)

type Command struct {
	ApplicationCommand *discordgo.ApplicationCommand
	Handler            func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate)
}

var (
	Commands = map[string]Command{
		"Отправить репорт": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:         discordgo.MessageApplicationCommand,
				Name:         "Отправить репорт",
				DMPermission: new(bool),
			},
			Handler: reportMessageCommandHandler,
		},
		"update-guild": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Type:        discordgo.ChatApplicationCommand,
				Name:        "update-guild",
				Description: "Обновить настройки сервера",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "report-channel",
						Description: "Канал для репортов",
						ChannelTypes: []discordgo.ChannelType{
							discordgo.ChannelTypeGuildText,
						},
						Required: false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "resolved-report-channel",
						Description: "Канал для рассмотренных репортов",
						ChannelTypes: []discordgo.ChannelType{
							discordgo.ChannelTypeGuildText,
						},
						Required: false,
					},
				},
				DMPermission: new(bool),
			},
			Handler: func(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
				err := session.InteractionRespond(interactionCreate.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags: discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					log.Println("Error responding to interaction: ", err)
					return
				}

				server := db.Guild{
					ID: interactionCreate.GuildID,
				}

				for _, option := range interactionCreate.ApplicationCommandData().Options {
					switch option.Name {
					case "report-channel":
						server.ReportChannelID = option.ChannelValue(session).ID
					case "resolved-report-channel":
						server.ResoledReportChannelID = option.ChannelValue(session).ID
					}
				}

				err = db.UpdateGuild(server)
				if err != nil {
					log.Printf("Error updating guild: %v", err)
					return
				}

				_, err = session.FollowupMessageCreate(interactionCreate.Interaction, true, &discordgo.WebhookParams{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Настройки сервера обновлены",
							Description: "Настройки сервера были обновлены.",
							Color:       DefaultEmbedColor,
						},
					},
				})
				if err != nil {
					log.Println("Error creating followup message: ", err)
					return
				}
			},
		},
	}
)

func CreateApplicationCommands(session *discordgo.Session) {
	for index, value := range Commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", value.ApplicationCommand)
		if err != nil {
			log.Panicf("Error creating '%v' command: %v", value.ApplicationCommand.Name, err)
		}
		log.Printf("Successfully created '%v' command", cmd.Name)

		if command, exists := Commands[index]; exists {
			command.ApplicationCommand = cmd
			Commands[index] = command
		}
	}
}
func RemoveApplicationCommands(session *discordgo.Session) {
	commands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		log.Printf("Error getting application commands: %v", err)
		return
	}
	for _, command := range commands {
		err = session.ApplicationCommandDelete(session.State.User.ID, "", command.ID)
		if err != nil {
			log.Panicf("Error deleting '%v' command: %v", command.Name, err)
		}
		log.Printf("Successfully deleted '%v' command", command.Name)
	}
}

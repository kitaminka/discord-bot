package interactions

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
)

var (
	Reasons = []Reason{
		{
			Name:        "Оффтоп",
			Description: "Сообщение не относится к теме канала.",
		},
		{
			Name:        "Реклама",
			Description: "Сообщение содержит рекламу.",
		},
	}
	reasonChoices = getReasonChoices()
)

type Reason struct {
	Name        string
	Description string
}

func getReasonChoices() []*discordgo.ApplicationCommandOptionChoice {
	var choices []*discordgo.ApplicationCommandOptionChoice
	for i, reason := range Reasons {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  reason.Name,
			Value: strconv.Itoa(i),
		})
	}
	return choices
}

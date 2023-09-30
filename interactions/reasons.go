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
	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(Reasons))
	for i, reason := range Reasons {
		choices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  reason.Name,
			Value: strconv.Itoa(i),
		}
	}
	return choices
}

package interactions

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
)

var (
	Reasons = []Reason{
		{
			Name:        "Оффтоп",
			Description: "**Вы отправляете сообщения, не соответствующие по тематике каналу!** Тематика канала указана в названии и описании канала.",
		},
		{
			Name:        "Неадекватное содержание",
			Description: "**Вы отправляете сообщения, не имеющие смысла.** Просьба общаться адекватно и уважительно взаимодействовать с участниками сервера.",
		},
		{
			Name:        "Флуд",
			Description: "**Вы отправляете неоднократно повторяющиеся сообщения.** Просьба соблюдать правила сервера.",
		},
		{
			Name:        "Язык вражды",
			Description: "**Вы общаетесь агрессивно или оскорбительно.** Просьба взаимодействовать с участниками сервера уважительно.",
		},
		{
			Name:        "Возраст",
			Description: "**Согласно условиям использования Discord, вам не может быть меньше 13 лет.** Просьба не упоминать это в своих сообщениях.",
		},
		{
			Name:        "Реклама",
			Description: "**Вы отправляете сообщения, содержащие рекламу.** Просьба не рекламировать сторонние ресурсы.",
		},
		{
			Name:        "Реклама Discord сервера",
			Description: "**Вы рекламируете своё Discord-сообщество.** Канал <#1167868890374738090> подходит только для рекламы внутриигровых сообществ.",
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

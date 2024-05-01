package msg

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type StructuredText struct {
	Text   string
	Fields []*StructuredTextField
}
type StructuredTextField struct {
	Emoji discordgo.Emoji
	Name  string
	Value string
}

func (structuredText StructuredText) ToString() string {
	result := structuredText.Text

	if len(structuredText.Fields) > 0 {
		result += "\n\n"
	}

	for _, field := range structuredText.Fields {
		messageEmoji := field.Emoji.MessageFormat()

		if len(messageEmoji) != 0 {
			result += fmt.Sprintf("%v **%v**: %v\n", messageEmoji, field.Name, field.Value)
		} else {
			result += fmt.Sprintf("**%v**: %v\n", field.Name, field.Value)
		}
	}

	return result
}

package msg

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type StructuredDescription struct {
	Text   string
	Fields []*StructuredDescriptionField
}
type StructuredDescriptionField struct {
	Emoji discordgo.Emoji
	Name  string
	Value string
}

func (structuredDescription StructuredDescription) ToString() string {
	result := structuredDescription.Text

	if len(structuredDescription.Fields) > 0 {
		result += "\n\n"
	}

	for _, field := range structuredDescription.Fields {
		messageEmoji := field.Emoji.MessageFormat()

		if len(messageEmoji) != 0 {
			result += fmt.Sprintf("%v **%v**: %v\n", messageEmoji, field.Name, field.Value)
		} else {
			result += fmt.Sprintf("**%v**: %v\n", field.Name, field.Value)
		}
	}

	return result
}

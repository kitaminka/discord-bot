package msg

import "fmt"

type StructuredDescription struct {
	Text   string
	Fields []*StructuredDescriptionField
}
type StructuredDescriptionField struct {
	Emoji Emoji
	Name  string
	Value string
}

func (structuredDescription StructuredDescription) ToString() string {
	result := structuredDescription.Text

	if len(structuredDescription.Fields) > 0 {
		result += "\n\n"
	}

	for _, field := range structuredDescription.Fields {
		if field.Emoji != "" {
			result += fmt.Sprintf("%v **%v**: %v\n", field.Emoji, field.Name, field.Value)
		} else {
			result += fmt.Sprintf("**%v**: %v\n", field.Name, field.Value)
		}
	}

	return result
}

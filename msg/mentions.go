package msg

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func UserMention(user *discordgo.User) string {
	return fmt.Sprintf("**%v** (%v)", user.Username, user.Mention())
}

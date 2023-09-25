package interactions

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

const (
	MuteDuration       = 36 * time.Hour
	ExtendedMutePeriod = 48 * time.Hour
	MuteWarningsCount  = 3
)

func isUserModerator(session *discordgo.Session, interaction *discordgo.Interaction, discordUser *discordgo.User) (bool, error) {
	perms, err := session.UserChannelPermissions(discordUser.ID, interaction.ChannelID)
	if err != nil {
		return false, err
	}

	guild, err := session.Guild(interaction.GuildID)
	if err != nil {
		return false, err
	}

	return perms&discordgo.PermissionModerateMembers != 0 || discordUser.ID == guild.OwnerID, nil
}

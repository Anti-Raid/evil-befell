package ext_types

import (
	"github.com/anti-raid/evil-befall/types/discordgo"
)

type GuildChannelWithPermissions struct {
	User    Permissions        `json:"user" description:"The permissions the user has in the channel"`
	Bot     Permissions        `json:"bot" description:"The permissions the bot has in the channel"`
	Channel *discordgo.Channel `json:"channel" description:"The channel object"`
}

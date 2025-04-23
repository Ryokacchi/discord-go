package utils

import "github.com/bwmarrin/discordgo"

func TotalMemberCount(guilds []*discordgo.Guild) (base int) {
	for _, v := range guilds {
		base += v.MemberCount
	}

	return base
}

package modules

import (
	"discord-go/commands"

	"github.com/bwmarrin/discordgo"
)

var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping":     commands.Ping,
	"botbilgi": commands.BotInfo,
}

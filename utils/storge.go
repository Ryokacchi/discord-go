package utils

import (
	"github.com/bwmarrin/discordgo"
)

var ApplicationCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Bot'un sağlıklı çalıştığını doğrulamak için bu komutu kullanabilirsiniz.",
	},
	{
		Name:        "botbilgi",
		Description: "Bot hakkında daha fazla bilgi edinmek için bu komutu kullanabilirsiniz.",
	},
}

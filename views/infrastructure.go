package views

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Color int = 0x5865f2
)

func Loading(i *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: Color,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%s (@%s) — Komut Yürütülüyor", i.Member.User.GlobalName, i.Member.User.Username),
			IconURL: i.Member.User.AvatarURL("128"),
		},
		Description: "Komutun cevap vermesi uzun sürebilir, anlayışınız için teşekkürler.",
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Komut yürütülüyor, lütfen bekleyiniz.",
			IconURL: "https://cdn.discordapp.com/emojis/1216348767096143962.gif",
		},
	}
}

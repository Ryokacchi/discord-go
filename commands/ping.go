package commands

import (
	"discord-go/db"
	"discord-go/views"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	database, websocket := db.Ping(), s.HeartbeatLatency()

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Color: views.Color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    fmt.Sprintf("%s (@%s) — WebSocket Gecikmesi", i.Member.User.GlobalName, i.Member.User.Username),
					IconURL: i.Member.User.AvatarURL("128"),
				},
				Description: "Burada bulunan gecikme değerleri anlık değildir; belirli zaman döngülerinde bir kez alınır ve kaydedilir.",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "WebSocket",
						Value:  fmt.Sprintf("```ansi\n\u001b[2;34m\u001b[1;34m%dms\u001b[0m\u001b[2;34m\u001b[0m\n```", websocket.Milliseconds()),
						Inline: true,
					},
					{
						Name:   "Veritabanı",
						Value:  fmt.Sprintf("```ansi\n\u001b[2;34m\u001b[1;34m%dms\u001b[0m\u001b[2;34m\u001b[0m\n```", database.Milliseconds()),
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("%s %d, best discord bot.", s.State.User.Username, time.Now().Year()),
				},
			},
		},
	},
	)
}

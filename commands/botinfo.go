package commands

import (
	"discord-go/db"
	"discord-go/utils"
	"discord-go/views"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

func BotInfo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	dbPing := db.Ping()
	latency := s.HeartbeatLatency()

	memoryStats, CPUUsage := utils.ReadMemoryStats(), utils.ReadCPUUsage()

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Color: views.Color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    fmt.Sprintf("%s (@%s) — Bot Bilgi", i.Member.User.GlobalName, i.Member.User.Username),
					IconURL: i.Member.User.AvatarURL("128"),
				},
				Description: fmt.Sprintf("%s verileri anlık olarak güncelleme gösterebilir ve tüm kullanıcı gerektiren veriler, kullanıcı kimliğini açığa çıkarmaz.", s.State.User.Username),
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "**• Genel Veriler:**",
						Value:  strings.Join([]string{fmt.Sprintf("`-` Sunucu Sayısı: `%d`", len(s.State.Guilds)), fmt.Sprintf("`-` Kullanıcı Sayısı: `%d`", utils.TotalMemberCount(s.State.Guilds))}, "\n"),
						Inline: true,
					},
					{
						Name:   "** **",
						Value:  strings.Join([]string{fmt.Sprintf("`-` WebSocket: `%dms`", latency.Milliseconds()), fmt.Sprintf("`-` Veritabanı: `%dms`", dbPing.Milliseconds())}, "\n"),
						Inline: true,
					},
					{
						Name:   "** **",
						Value:  "** **",
						Inline: true,
					},
					{
						Name:   "**• Sayısal Veriler:**",
						Value:  strings.Join([]string{fmt.Sprintf("`-` Komut Sayısı: `%d`", len(utils.ApplicationCommands)), fmt.Sprintf("`-` Etkinlik Sayısı: `%d`", utils.ActiveHandlers)}, "\n"),
						Inline: true,
					},
					{
						Name:   "**  **",
						Value:  strings.Join([]string{fmt.Sprintf("`-` Kullanılan Bellek: `%s (%s)`", utils.FormatBytes(memoryStats.Runtime), utils.FormatBytes(memoryStats.System.Available)), fmt.Sprintf("`-` İşlemci Kullanımı: `%.2f%%`", CPUUsage)}, "\n"),
						Inline: true,
					},
					{
						Name:   "** **",
						Value:  "** **",
						Inline: true,
					},
					{
						Name:   "• Tarihsel veriler:",
						Value:  strings.Join([]string{fmt.Sprintf("`-` Bot başlama süresi: **%s**", utils.FormatTime(humanize.Time(utils.Uptime))), ""}, "\n"),
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

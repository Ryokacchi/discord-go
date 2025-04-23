package events

import (
	"discord-go/db"
	"discord-go/utils"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		commandData := i.ApplicationCommandData()

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		if err != nil {
			log.Printf("Error responding to command: %v", err)
			return
		}

		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: utils.StringPtr("komut yükleniyor"),
		})
		if err != nil {
			log.Printf("Error responding to command: %v", err)
		}

		if commandData.Name == "ping" {
			dbPing := db.Ping()
			latency := s.HeartbeatLatency()

			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Color:       0x5865F2,
						Title:       "WebSocket Gecikmesi",
						Description: "Burada bulunan gecikme değerleri anlık değildir; belirli zaman döngülerinde bir kez alınır ve kaydedilir.",
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "WebSocket",
								Value: fmt.Sprintf("```ansi\n\u001b[2;34m\u001b[1;34m%dms\u001b[0m\u001b[2;34m\u001b[0m\n```", latency.Milliseconds()),
							},
							{
								Name:  "Veritabanı",
								Value: fmt.Sprintf("```ansi\n\u001b[2;34m\u001b[1;34m%dms\u001b[0m\u001b[2;34m\u001b[0m\n```", dbPing.Milliseconds()),
							},
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text: fmt.Sprintf("%s %d, best discord bot.", s.State.User.Username, time.Now().Year()),
						},
					},
				},
			},
			)

			if err != nil {
				log.Printf("Error responding to command: %v", err)
			}
		}

		if commandData.Name == "botbilgi" {
			dbPing := db.Ping()
			latency := s.HeartbeatLatency()

			memoryStats, CPUUsage := utils.ReadMemoryStats(), utils.ReadCPUUsage()

			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Color:       0x5865F2,
						Title:       "Bot Bilgi",
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

			if err != nil {
				log.Printf("Error responding to command: %v", err)
			}
		}
	}
}

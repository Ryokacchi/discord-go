package main

import (
	database "discord-go/Database"
	functions "discord-go/Functions"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/mem"
)

var (
	activeHandlers int       = 0
	uptime         time.Time = time.Now()
)

func TrackHandler(handler interface{}) interface{} {
	activeHandlers++
	return handler
}

func translateToTurkish(en string) string {
	replacements := map[string]string{
		"years ago":   "yıl önce",
		"year ago":    "yıl önce",
		"months ago":  "ay önce",
		"month ago":   "ay önce",
		"days ago":    "gün önce",
		"day ago":     "gün önce",
		"hours ago":   "saat önce",
		"hour ago":    "saat önce",
		"minutes ago": "dakika önce",
		"minute ago":  "dakika önce",
		"seconds ago": "saniye önce",
		"second ago":  "saniye önce",
		" ago":        " önce",
	}

	for k, v := range replacements {
		if strings.Contains(en, k) {
			return strings.Replace(en, k, v, 1)
		}
	}

	return en
}

/*
*
* exit blocks the main goroutine until an interrupt or termination signal is received.
* Once a signal is received, it prints a shutdown message and allows the program to exit gracefully.
 */
func exit() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	fmt.Println("Shutdown signal received, shutting down the application...")
}

func initialize(token string) *discordgo.Session {
	dg, err := discordgo.New(strings.Join([]string{"Bot", token}, " "))
	if err != nil {
		log.Fatalf("Failed to create Discord bot: %v", err)
	}

	/** Guilds Intents */
	dg.Identify.Intents = discordgo.IntentGuilds

	return dg
}

func StringPtr(s string) *string {
	return &s
}

/** Ready Event */
func ready(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println("Bot successfully connected to Discord.")
	fmt.Println("Logged in as " + r.User.Username + "!")
}

/** InteractionCreate Event */
func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			Content: StringPtr("komut yükleniyor"),
		})
		if err != nil {
			log.Printf("Error responding to command: %v", err)
		}

		if commandData.Name == "ping" {
			dbPing := database.Ping()
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
			dbPing := database.Ping()
			latency := s.HeartbeatLatency()

			vm, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("Error responding to memory: %v", err)
				return
			}

			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			cpuUsage, err := functions.GetCPUUsage()
			if err != nil {
				log.Printf("Error responding to cpu: %v", err)
				return
			}

			commands, err := s.ApplicationCommands(s.State.User.ID, "")
			if err != nil {
				fmt.Printf("Failed to fetch application commands: %v\n", err)
				return
			}

			userCount := 0
			for _, guild := range s.State.Guilds {
				userCount += guild.MemberCount
			}

			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Color:       0x5865F2,
						Title:       "Bot Bilgi",
						Description: fmt.Sprintf("%s verileri anlık olarak güncelleme gösterebilir ve tüm kullanıcı gerektiren veriler, kullanıcı kimliğini açığa çıkarmaz.", s.State.User.Username),
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "**• Genel Veriler:**",
								Value:  strings.Join([]string{fmt.Sprintf("`-` Sunucu Sayısı: `%d`", len(s.State.Guilds)), fmt.Sprintf("`-` Kullanıcı Sayısı: `%d`", userCount)}, "\n"),
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
								Value:  strings.Join([]string{fmt.Sprintf("`-` Komut Sayısı: `%d`", len(commands)), fmt.Sprintf("`-` Etkinlik Sayısı: `%d`", activeHandlers)}, "\n"),
								Inline: true,
							},
							{
								Name:   "**  **",
								Value:  strings.Join([]string{fmt.Sprintf("`-` Kullanılan Bellek: `%s (%s)`", functions.FormatBytes(memStats.Alloc), functions.FormatBytes(vm.Available)), fmt.Sprintf("`-` İşlemci Kullanımı: `%.2f%%`", cpuUsage)}, "\n"),
								Inline: true,
							},
							{
								Name:   "** **",
								Value:  "** **",
								Inline: true,
							},
							{
								Name:   "• Tarihsel veriler:",
								Value:  strings.Join([]string{fmt.Sprintf("`-` Bot başlama süresi: **%s**", translateToTurkish(humanize.Time(uptime))), ""}, "\n"),
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

func main() {
	cfg := functions.ConfigLoader()
	dg := initialize(cfg.Section("bot").Key("token").String())

	/** Registers an event handler for the specified event. */
	dg.AddHandler(TrackHandler(ready))
	dg.AddHandler(TrackHandler(interactionCreate))

	/** Sets up the MongoDB connection and makes the client accessible globally. */
	database.Connect()

	/** Open the connection to the Discord bot */
	err := dg.Open()
	if err != nil {
		log.Fatalf("Failed to connect to Discord bot: %v", err)
	}

	_, err = dg.ApplicationCommandBulkOverwrite(dg.State.Application.ID, "", functions.ApplicationCommands)
	if err != nil {
		log.Fatalf("Failed to create application commands: %v", err)
	}

	exit()
}

package main

import (
	database "discord-go/Database"
	functions "discord-go/Functions"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

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

/** Ready Event */
func ready(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println("Bot successfully connected to Discord.")
	fmt.Println("Logged in as " + r.User.Username + "!")
}

/** InteractionCreate Event */
func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		commandData := i.ApplicationCommandData()

		if commandData.Name == "ping" {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			if err != nil {
				log.Printf("Error responding to command: %v", err)
				return
			}

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
	}
}

func main() {
	cfg := functions.ConfigLoader()
	dg := initialize(cfg.Section("bot").Key("token").String())

	/** Registers an event handler for the specified event. */
	dg.AddHandler(ready)
	dg.AddHandler(interactionCreate)

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

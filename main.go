package main

import (
	functions "discord-go/Functions"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/ini.v1"
)

var (
	configPath = "config.ini"
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

/*
*
* configLoader loads the configuration from config.ini file.
* If the config file cannot be loaded, it logs the error and exits the program.
 */
func configLoader() *ini.File {
	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config.ini file: %v", err)
	}

	return cfg
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
			latency := s.HeartbeatLatency()

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Color:       0x5865F2,
							Title:       "WebSocket Gecikmesi",
							Description: "Burada bulunan gecikme değerleri anlık değildir; belirli zaman döngülerinde bir kez alınır ve kaydedilir.",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "WebSocket",
									Value: fmt.Sprintf("```ansi\n\u001b[2;34m\u001b[1;34m%dms\u001b[0m\u001b[2;34m\u001b[0m\n```", latency.Milliseconds()),
								},
							},
							Footer: &discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("%s %d, best discord bot.", s.State.User.Username, time.Now().Year()),
							},
						},
					},
				},
			})

			if err != nil {
				log.Printf("Error responding to command: %v", err)
			}
		}
	}
}

func main() {
	cfg := configLoader()
	dg := initialize(cfg.Section("bot").Key("token").String())

	/** Registers an event handler for the specified event. */
	dg.AddHandler(ready)
	dg.AddHandler(interactionCreate)

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

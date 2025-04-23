package main

import (
	"discord-go/db"
	"discord-go/events"
	"discord-go/utils"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// exit blocks the main goroutine until an interrupt or termination signal is received.
// Once a signal is received, it prints a shutdown message and allows the program to exit gracefully.
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

	// Guilds Intents
	dg.Identify.Intents = discordgo.IntentGuilds

	return dg
}

func main() {
	cfg := utils.ConfigLoader()
	dg := initialize(cfg.Section("bot").Key("token").String())

	// Registers an event handler for the specified event.
	dg.AddHandler(utils.TrackHandler(events.Ready))
	dg.AddHandler(utils.TrackHandler(events.InteractionCreate))

	// Sets up the MongoDB connection and makes the client accessible globally.
	db.Connect()

	// Open the connection to the Discord bot
	err := dg.Open()
	if err != nil {
		log.Fatalf("Failed to connect to Discord bot: %v", err)
	}

	_, err = dg.ApplicationCommandBulkOverwrite(dg.State.Application.ID, "", utils.ApplicationCommands)
	if err != nil {
		log.Fatalf("Failed to create application commands: %v", err)
	}

	exit()
}

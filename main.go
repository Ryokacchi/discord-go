package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

func ready(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println("Bot successfully connected to Discord.")
	fmt.Println("Logged in as " + r.User.Username + "!")
}

func main() {
	cfg := configLoader()
	dg := initialize(cfg.Section("bot").Key("token").String())

	/** Registers the 'ready' event handler to handle the "ready" event from Discord. */
	dg.AddHandler(ready)

	/** Open the connection to the Discord bot */
	err := dg.Open()
	if err != nil {
		log.Fatalf("Failed to connect to Discord bot: %v", err)
	}

	exit()
}

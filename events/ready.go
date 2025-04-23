package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println("Bot successfully connected to Discord.")
	fmt.Println("Logged in as " + r.User.Username + "!")
}

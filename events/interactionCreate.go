package events

import (
	"discord-go/modules"
	"discord-go/utils"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
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

		if handler, ok := modules.CommandHandlers[commandData.Name]; ok {
			handler(s, i)
		} else {
			fmt.Printf("command not found: %s\n", commandData.Name)
		}
	}
}

package events

import (
	"discord-go/modules"
	"discord-go/views"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		commandData := i.ApplicationCommandData()

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{views.Loading(i)},
		})

		if handler, ok := modules.CommandHandlers[commandData.Name]; ok {
			handler(s, i)
		} else {
			fmt.Printf("command not found: %s\n", commandData.Name)
		}
	}
}

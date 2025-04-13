package functions

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Bot'un sağlıklı çalıştığını doğrulamak için bu komutu kullanabilirsiniz.",
	},
}

/*
* Registers all defined commands with Discord, creating them one by one.
 */
func PublishCommands(dg *discordgo.Session) {
	cmds := make([]*discordgo.ApplicationCommand, len(Commands))

	for i, v := range Commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", v)
		if err != nil {
			log.Fatalf("Failed to create application command: %v", err)
		}

		cmds[i] = cmd
		fmt.Println("Command '" + cmd.Name + "' registered successfully.")
	}
}

/*
* Retrieves and deletes all registered commands from Discord.
 */
func DeleteCommands(dg *discordgo.Session) {
	cmds, err := dg.ApplicationCommands(dg.State.User.ID, "")
	if err != nil {
		log.Fatalf("An error occurred while retrieving commands: %v", err)
	}

	for _, v := range cmds {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, "", v.ID)
		if err != nil {
			log.Fatalf("An error occurred while deleting the command: %v", err)
		}

		fmt.Println("Command '" + v.Name + "' was successfully deleted.")
	}
}

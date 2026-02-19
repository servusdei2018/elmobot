package commands

import "github.com/bwmarrin/discordgo"

var Ask = &discordgo.ApplicationCommand{
	Name:        "ask",
	Description: "Ask Elmo a question",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "question",
			Description: "Your question for Elmo",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
}

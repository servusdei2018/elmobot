package commands

import "github.com/bwmarrin/discordgo"

var Ping = &discordgo.ApplicationCommand{
	Name:        "ping",
	Description: "Play ping pong with Elmo",
}

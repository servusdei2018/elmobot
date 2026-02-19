package commands

import (
	"github.com/bwmarrin/discordgo"

	"github.com/servusdei2018/elmobot/pkg/handlers"
)

var (
	Cmds = []*discordgo.ApplicationCommand{
		Ping,
		Ask,
	}

	Handlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": handlers.Ping,
		"ask":  handlers.Ask,
	}
)

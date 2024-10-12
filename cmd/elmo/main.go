package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"

	"github.com/servusdei2018/elmobot/pkg/commands"
)

var (
	token = flag.String("token", "", "Discord bot token")
	s     *discordgo.Session
)

func init() {
	flag.Parse()

	var err error
	s, err = discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatalf("error connecting: %v", err)
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commands.Handlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("error opening the session: %v", err)
	}

	log.Println("adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands.Cmds))
	for i, v := range commands.Cmds {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("error creating command '%v': %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("press Ctrl+C to exit")
	<-stop

	registeredCommands, err = s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		log.Fatalf("error fetching registered commands: %v", err)
	}

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("error deleting '%v' command: %v", v.Name, err)
		}
	}

	log.Println("shutting down")
}

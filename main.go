package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"trup/command"
)

var (
	prefix = "."
	env    = command.Env{
		RoleMod:         os.Getenv("ROLE_MOD"),
		ChannelShowcase: os.Getenv("CHANNEL_SHOWCASE"),
	}
	botId string
)

func main() {
	var (
		token = os.Getenv("TOKEN")
	)

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Failed on discordgo.New(): %s\n", err)
	}

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		botId = r.User.ID
		s.UpdateStatus(0, ".help")
	})
	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		log.Fatalf("Failed on discord.Open(): %s\n", err)
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in messageCreate. r: %#v; Message(%s): %s;\n", r, m.ID, m.Content)
		}
	}()

	if m.ChannelID == env.ChannelShowcase {
		for _, a := range m.Attachments {
			if a.Width > 0 {
				command.UpdateSysinfoImage(m.Author.ID, a.URL)
				break
			}
		}

		err := s.MessageReactionAdd(m.ChannelID, m.ID, "❤")
		if err != nil {
			log.Printf("Error on adding reaction ❤ to new showcase message(%s): %s\n", m.ID, err)
			return
		}
	}

	var mentionsBot bool
	for _, m := range m.Mentions {
		if m.ID == botId {
			mentionsBot = true
			break
		}
	}
	if mentionsBot {
		s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" need help? Type `.help`")
		return
	}

	if strings.HasPrefix(m.Content, prefix) {
		args := strings.Fields(m.Content[len(prefix):])
		context := command.Context{
			Env:     &env,
			Session: s,
			Message: m.Message,
		}

		if len(args) == 0 {
			return
		}

		if args[0] == "help" {
			command.Help(&context, args)
			return
		}

		var found bool
		allKeys := make([]string, 0, len(command.Commands))
		for name, cmd := range command.Commands {
			allKeys = append(allKeys, name)
			if !found && args[0] == name {
				found = true
				cmd.Exec(&context, args)
			}
		}
		if !found {
			// this will need to be either disabled
			// or need a workaround for situations like "..." when PREFIX=.
			//s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" command not found. Available commands: "+strings.Join(allKeys, ", "))
		}
	}
}

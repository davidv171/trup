package command

import "github.com/bwmarrin/discordgo"

type Env struct {
	RoleMod         string
	ChannelShowcase string
}

type Context struct {
	Env     *Env
	Session *discordgo.Session
	Message *discordgo.Message
}

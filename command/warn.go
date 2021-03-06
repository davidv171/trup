package command

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"log"
	"strings"
	"trup/db"
)

const warnUsage = "warn <@user> <reason>"

func warn(ctx *Context, args []string) {
	if len(args) < 3 {
		ctx.Reply("not enough arguments.")
		return
	}

	var (
		user   = parseMention(args[1])
		reason = strings.Join(args[2:], " ")
	)

	w := db.NewWarn(ctx.Message.Author.ID, user, reason)
	err := w.Save()
	if err != nil {
		msg := fmt.Sprintf("Failed to save your warning. Error: %s", err)
		log.Println(msg)
		ctx.Reply(msg)
		return
	}

	var nth string
	warnCount, err := db.CountWarns(user)
	if err != nil {
		log.Printf("Failed to count warns for user %s; Error: %s\n", user, err)
	}
	if warnCount > 0 {
		nth = " for the " + humanize.Ordinal(warnCount) + " time"
	}

	err = db.NewNote(ctx.Message.Author.ID, user, "User was warned for: "+reason).Save()
	if err != nil {
		log.Printf("Failed to save warning note. Error: %s\n", err)
	}

	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("<@%s> Has been warned%s with reason: %s.", user, nth, reason))
}

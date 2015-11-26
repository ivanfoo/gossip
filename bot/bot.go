package bot

import (
	"log"
	"os"
	"os/user"
	"strings"

	_ "github.com/ivanfoo/rtop-bot/utils"
	"github.com/nlopes/slack"
)

type BotOptions struct {
	Username   string
	SSHKeyPath string
	SlackToken string
}

type Bot struct {
	botOptions BotOptions
	SlackUserID string
	SystemUser *user.User
}

func NewBot(opts BotOptions) *Bot {
	b := new(Bot)
	b.botOptions = opts

	if opts.Username != "" {
		b.SystemUser, _ = user.Lookup(opts.Username)
	} else {
		b.SystemUser, _ = user.Current()
	}

	if opts.SSHKeyPath == "" {
		b.botOptions.SSHKeyPath = b.SystemUser.HomeDir
	}

	return b
}

func (b *Bot) DoSlack() {
	api := slack.New(b.botOptions.SlackToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {

			case *slack.ConnectedEvent:
				b.SlackUserID = ev.Info.User.ID
				log.Println("Infos:", ev.Info)
				log.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				if b.botBeingAsked(ev.Msg.Text) {
					rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", ev.Msg.Channel))	
				}

			case *slack.InvalidAuthEvent:
				log.Printf("Invalid credentials")
				os.Exit(1)
			}	
		}
	}
}

func (b *Bot) botBeingAsked(slackMessage string) bool {
	botMention := "<@" + b.SlackUserID
	return strings.HasPrefix(slackMessage, botMention)
}


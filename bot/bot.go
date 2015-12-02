package bot

import (
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/ivanfoo/gossip/commands"
	"github.com/ivanfoo/gossip/utils"
	"github.com/nlopes/slack"
)

type BotOptions struct {
	Username   string
	SSHKeyPath string
	SlackToken string
}

type Bot struct {
	botOptions  BotOptions
	SlackUserID string
	SystemUser  *user.User
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
				log.Println(ev.Msg.Text)
				if b.botBeingAsked(ev.Msg.Text) {
					response := b.runCommand(b.parseMessage(ev.Msg.Text))
					rtm.SendMessage(rtm.NewOutgoingMessage(response, ev.Msg.Channel))
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

func (b *Bot) parseMessage(slackMessage string) (action string, target string) {
	parts := strings.Fields(slackMessage)
	parts[2] = utils.CleanHostname(parts[2])
	return parts[1], parts[2] 
}

func (b *Bot) runCommand(action string, target string) string {
	sshClient := utils.SSHConnect(b.SystemUser.Username, target, b.botOptions.SSHKeyPath)
	defer sshClient.Close()
	c := commands.NewCommand(sshClient)
	output := c.Run()
	return output
}
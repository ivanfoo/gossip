package bot

import (
	"fmt"
	"log"
	"os/user"
	"strings"

	"github.com/ivanfoo/gossip/monitors"
	"github.com/ivanfoo/gossip/utils"

	"golang.org/x/crypto/ssh"
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
		b.botOptions.SSHKeyPath = b.SystemUser.HomeDir + "/.ssh/" + "id_rsa" 
	}

	return b
}

func (b *Bot) Chat() {
    ws, id := utils.SlackConnect(b.botOptions.SlackToken)
    fmt.Println("bot ready, ^C exits")

    for {
        m, err := utils.GetMessage(ws)
        if err != nil {
            log.Fatal(err)
        }

        if m.Type == "message" && strings.HasPrefix(m.Text, "<@" + id) {
	        go func(m utils.Message) {
	        	action, target := b.parseMessage(m.Text)
	            	m.Text = b.runMonitor(action, target)
	            	utils.PostMessage(ws, m)
	        }(m)
        }
    }
}

func (b *Bot) botMentioned(slackMessage string) bool {
	botMention := "<@" + b.SlackUserID
	return strings.HasPrefix(slackMessage, botMention)
}

func (b *Bot) parseMessage(slackMessage string) (action string, target string) {
	parts := strings.Fields(slackMessage)
	parts[2] = utils.CleanHostname(parts[2])
	return parts[1], parts[2] 
}

func (b *Bot) runMonitor(action string, target string) string {
	sshClient := utils.SSHConnect(b.SystemUser.Username, target, b.botOptions.SSHKeyPath)
	defer sshClient.Close()

	var actions = map[string]func(client *ssh.Client) string {
	"containers": monitors.RunContainersMonitor,
	"uptime": monitors.RunUptimeMonitor,
	}

	return (actions[action](sshClient))
	//fmt.Println(fn())
}

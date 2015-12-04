package bot

import (
	"os/user"
	"strings"

	"github.com/ivanfoo/gossip/monitors"
	"github.com/ivanfoo/gossip/utils"
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
		b.botOptions.SSHKeyPath = b.SystemUser.HomeDir + ".ssh/" + "id_rsa" 
	}

	return b
}

func (b *Bot) Run() {
	doSlack(b)
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
	m := monitors.NewStatsMonitor(sshClient)
	output := m.Run()
	return output
}
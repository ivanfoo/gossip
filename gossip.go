package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ivanfoo/gossip/bot"
)

const (
	VERSION = "1.0.0"
)

func usage() {
	fmt.Printf(	
		`	
Usage:
    gossip -s slackBotToken [-u user] [-i identity_file]

where:
    -s slackBotToken is the API token for the Slack bot
    -u user
    -i identity_file
`, VERSION)
	os.Exit(1)
}

func main() {
	var SlackToken = flag.String("s", "", "slack bot token")
	var Username = flag.String("u", "", "ssh user to use")
	var SSHKeyPath = flag.String("i", "", "private key to use")

	flag.Parse()

	if *SlackToken == "" {
		usage()
	}

	bot := bot.NewBot(bot.BotOptions{
		Username:   *Username,
		SSHKeyPath: *SSHKeyPath,
		SlackToken: *SlackToken,
	})

	bot.Chat()
}

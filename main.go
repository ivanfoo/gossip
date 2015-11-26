package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ivanfoo/rtop-bot/bot"
)

const (
	VERSION = "0.1"
)

func usage() {
	fmt.Printf(	
		`	
Usage:
    rtop-bot -s slackBotToken [-u user] [-i identity_file]

where:
    -s slackBotToken is the API token for the Slack bot
    -u user
    -i identity_file
`, VERSION)
	os.Exit(1)
}

func main() {
	var SlackToken = flag.String("s", "", "create Slack bot")
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

	bot.DoSlack()
}

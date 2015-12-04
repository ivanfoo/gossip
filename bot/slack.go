package bot	

import (
	"log"
	"os"
	"github.com/nlopes/slack"
)

func doSlack(b *Bot) {
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
				if b.botMentioned(ev.Msg.Text) {
					response := b.runMonitor(b.parseMessage(ev.Msg.Text))
					rtm.SendMessage(rtm.NewOutgoingMessage(response, ev.Msg.Channel))
				}

			case *slack.InvalidAuthEvent:
				log.Printf("Invalid credentials")
				os.Exit(1)
			}
		}
	}
}
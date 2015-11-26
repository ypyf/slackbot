package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"os"
	"regexp"
	"strings"

	"slackbot/plugin"
)

var debug_mode bool = true

const (
	botID = "U0F8EH14Y"
	botDMChannelID = "D0F8EH15W"
	botChannelID = "C0F8FQ2GH"
	mxTeamID = "T0F8CMAU9"
)

func checkMessage(msg string) (string, bool) {
	r := regexp.MustCompile("^<@([\\d\\w]+)>:(.*)")
	s := r.FindStringSubmatch(msg)
	if len(s) == 0 {
		return "", false
	}
	return s[2], s[1] == botID
}

func handleCommand(rtm *slack.RTM, botName, channel, sender, text string) {
	for _, v := range plugin.BotCommands {
		if v.Matches(text) {
			msg := plugin.NewMessage(rtm, botName, text, channel)
			if err := v.Respond(msg); err != nil {
				var reply string
				if debug_mode {
					reply = fmt.Sprintf("Opps! %s遇到了点麻烦:\n%s", botName, err.Error())
				} else {
					reply = fmt.Sprintf("Opps! %s遇到了点麻烦，正在紧张处理中...", botName)
				}
				rtm.SendMessage(rtm.NewOutgoingMessage(reply, channel))
			}
			break
		}
	}
	return
}

func handleMessage(rtm *slack.RTM) {
	// 注册命令处理器
	plugin.BotCommands = append(plugin.BotCommands, new(plugin.Help))
	plugin.BotCommands = append(plugin.BotCommands, new(plugin.Hello))
	plugin.BotCommands = append(plugin.BotCommands, new(plugin.Time))
	plugin.BotCommands = append(plugin.BotCommands, new(plugin.Service))
	plugin.BotCommands = append(plugin.BotCommands, new(plugin.Nginx))
	plugin.BotCommands = append(plugin.BotCommands, new(plugin.Joke))

	user, err := rtm.GetUserInfo(botID)
	if err != nil {
		fmt.Println(err)
		return
	}
	botName := user.Profile.FirstName + " " + user.Profile.LastName

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch evt := msg.Data.(type) {
			case *slack.HelloEvent:
			// Ignore hello
			case *slack.ConnectedEvent:
				fmt.Println("Info:", evt.Info)
				fmt.Println("Connection counter:", evt.ConnectionCount)
				rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))
			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", evt)
				if evt.Channel == botChannelID && evt.Team == mxTeamID {
					text, sendToMXBot := checkMessage(evt.Text)
					if sendToMXBot && len(text) > 0 {
						//fmt.Printf("接收到机器人指令： %v", evt)
						go handleCommand(rtm, botName, evt.Channel, evt.User, strings.TrimSpace(text))
					}
				} else if evt.Channel == botDMChannelID && evt.Team == mxTeamID {
					go handleCommand(rtm, botName, evt.Channel, evt.User, strings.TrimSpace(evt.Text))
				}
			case *slack.ChannelJoinedEvent:
			// Ignore
			case *slack.PresenceChangeEvent:
				fmt.Printf("Presence Change: %v\n", evt)
			case *slack.LatencyReport:
				fmt.Printf("Current latency: %v\n", evt.Value)
			case *slack.RTMError:
				fmt.Printf("Error: %s\n", evt.Error())
			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				return
			default:
			// Ignore other events...
			}
		}
	}
}

func main() {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	api.SetDebug(true)
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	handleMessage(rtm)
}

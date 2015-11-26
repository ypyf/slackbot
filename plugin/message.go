package plugin

import (
	"github.com/nlopes/slack"
)

type Commander interface {
	Matches(text string) bool
	Respond(msg *Message) error
	Help() string
}

var BotCommands []Commander

type Message struct {
	*slack.RTM
	BotName string
	Text    string
	Channel string
	buffer  string
}

func NewMessage(rtm *slack.RTM, botName, text, channel string) *Message {
	return &Message{rtm, botName, text, channel, ""}
}

func (m *Message) Send(msg string) {
	m.buffer += msg + "\n"
}

func (m *Message) Done() {
	m.SendMessage(m.NewOutgoingMessage(m.buffer, m.Channel))
}



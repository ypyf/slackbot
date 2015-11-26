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
	Text    string
	Channel string
	buffer  string
}

func NewMessage(rtm *slack.RTM, text, channel string) *Message {
	return &Message{rtm, text, channel, ""}
}

func (m *Message) Send(msg string) {
	m.buffer += msg + "\n"
}

func (m *Message) Done() {
	m.SendMessage(m.NewOutgoingMessage(m.buffer, m.Channel))
}



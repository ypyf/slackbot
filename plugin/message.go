package plugin

import (
	"github.com/nlopes/slack"
)

const (
	S_INIT = iota
	S_PROCESSING = iota
)

type Commander interface {
	Matches(text string) bool
	Respond(msg *Message) error
	Help() string
}

var BotCommands []Commander

type Session struct {
	Handler Commander
	Status int32
	params interface{}
}

func (s *Session) ResetSession() {
	s.Status = S_INIT
}


type Message struct {
	*slack.RTM
	*Session
	BotName string
	Text    string
	Channel string
	buffer  string
}

func NewMessage(rtm *slack.RTM, session *Session, botName, text, channel string) *Message {
	return &Message{rtm, session, botName, text, channel, ""}
}

func (m *Message) Send(msg string) {
	m.buffer += msg + "\n"
}

func (m *Message) Done() {
	m.SendMessage(m.NewOutgoingMessage(m.buffer, m.Channel))
}



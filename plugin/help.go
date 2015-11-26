package plugin

import (
	"strings"
)

type Help struct {
	Sender string // sender id
}

func (self *Help) Matches(text string) bool {
	return strings.HasPrefix(text, "help") || strings.Contains(text, "帮助")
}

func (self *Help) Respond(msg *Message) error {
	tokens := strings.Fields(msg.Text)
	if len(tokens) >= 2 {
		for _, v := range BotCommands {
			if v.Matches(tokens[1]) {
				msg.Send(v.Help())
				break
			}
		}
	} else {
		msg.Send("当前可用的命令:")
		for _, v := range BotCommands {
			tokens := strings.Fields(v.Help())
			if len(tokens) > 0 {
				msg.Send(tokens[0])
			}
		}
	}
	msg.Done()
	return nil
}

func (self *Help) Help() string {
	return "help - 你可以尝试输入 help <cmd> 来获得各种命令的帮助信息."
}
